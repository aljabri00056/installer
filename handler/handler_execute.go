package handler

import (
	"errors"
	"strings"
	"time"

	"github.com/aljabri00056/installer/handler/provider"
	"github.com/aljabri00056/installer/logger"
)

func (h *Handler) execute(provider provider.Provider, q Query) (Result, error) {
	key := q.cacheKey()
	h.cacheMut.Lock()
	if h.cache == nil {
		h.cache = map[string]Result{}
	}
	cached, ok := h.cache[key]
	h.cacheMut.Unlock()
	if ok && time.Since(cached.Timestamp) < cacheTTL {
		return cached, nil
	}
	ts := time.Now()

	release, assets, err := h.getAssets(provider, q)

	if err != nil {
		return Result{}, err
	}
	if q.Release == "" && release != "" {
		logger.Debug("detected release: %s", release)
		q.Release = release
	}
	hasM1Asset := false
	for _, a := range assets {
		if a.IsMacM1() {
			hasM1Asset = true
			break
		}
	}
	result := Result{
		Timestamp: ts,
		Query:     q,
		Assets:    assets,
		Version:   release,
		M1Asset:   hasM1Asset,
	}

	h.cacheMut.Lock()
	h.cache[key] = result
	h.cacheMut.Unlock()
	return result, nil
}

func (h *Handler) getAssets(_provider provider.Provider, q Query) (string, []provider.Asset, error) {
	user := q.User
	repo := q.Program
	release := q.Release

	logger.Debug("fetching asset info for %s/%s@%s", user, repo, release)

	version, assets, err := _provider.GetReleaseAssets(user, repo, release, q.Token)
	if err != nil {
		return "", nil, err
	}

	if len(assets) == 0 {
		return version, nil, errors.New("no assets found")
	}

	filtered := []provider.Asset{}
	index := map[string]bool{}

	for _, asset := range assets {
		if q.Token != "" && q.Private {
			asset.DownloadURL = asset.URL
		}
		fext := getFileExt(asset.Name)
		if fext == "" && asset.Size > 1024*1024 {
			fext = ".bin" // +1MB binary
		}

		switch fext {
		case ".bin", ".zip", ".tar.bz", ".tar.bz2", ".bz2", ".gz", ".tar.gz", ".tgz", ".tar.xz":
		default:
			logger.Debug("fetched asset has unsupported file type: %s (ext '%s')", asset.Name, fext)
			continue
		}

		if q.Include != "" {
			skip := true
			includes := strings.Split(q.Include, ",")
			for _, include := range includes {
				if strings.Contains(asset.Name, include) {
					skip = false
				}
			}
			if skip {
				continue
			}
		}

		os := getOS(asset.Name)
		arch := getArch(asset.Name)

		if os == "" {
			logger.Debug("fetched asset has unknown os: %s", asset.Name)
			continue
		}
		if arch == "" {
			continue
		}

		logger.Debug("fetched asset: %s", asset.Name)

		asset.OS = os
		asset.Arch = arch
		asset.Type = fext

		if !index[asset.Key()] {
			index[asset.Key()] = true
			filtered = append(filtered, asset)
		}
	}

	if len(filtered) == 0 {
		return version, nil, errors.New("no downloads found for this release")
	}

	return version, filtered, nil
}
