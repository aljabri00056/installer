package provider

import (
	"fmt"
)

type GitHub struct {
	BaseProvider
	BaseURL string
}

type ghAsset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Size               int    `json:"size"`
	URL                string `json:"url"`
}

type ghRelease struct {
	Assets  []ghAsset `json:"assets"`
	Name    string    `json:"name"`
	TagName string    `json:"tag_name"`
	URL     string    `json:"url"`
}

type ghRepo struct {
	Private bool `json:"private"`
}

func (g *GitHub) GetRepo(user, repo, token string) (*RepoInfo, error) {
	url := fmt.Sprintf(g.BaseURL+"/repos/%s/%s", user, repo)
	var res ghRepo
	if err := g.get(url, token, &res); err != nil {
		return nil, err
	}
	return &RepoInfo{Private: res.Private}, nil
}

func (g *GitHub) GetReleaseAssets(user, repo, release, token string) (string, []Asset, error) {
	var assets []Asset
	var version string

	url := fmt.Sprintf(g.BaseURL+"/repos/%s/%s/releases", user, repo)

	if release == "" || release == "latest" {
		url += "/latest"
		var resp ghRelease
		if err := g.get(url, token, &resp); err != nil {
			return "", nil, err
		}
		version = resp.TagName
		for _, a := range resp.Assets {
			assets = append(assets, Asset{
				Name:        a.Name,
				URL:         a.URL,
				Size:        a.Size,
				DownloadURL: a.BrowserDownloadURL,
			})
		}
	} else {
		version = release
		url = fmt.Sprintf(g.BaseURL+"/repos/%s/%s/releases/tags/%s", user, repo, release)
		var resp ghRelease
		if err := g.get(url, token, &resp); err != nil {
			return "", nil, err
		}
		for _, a := range resp.Assets {
			assets = append(assets, Asset{
				Name:        a.Name,
				URL:         a.URL,
				Size:        a.Size,
				DownloadURL: a.BrowserDownloadURL,
			})
		}
	}

	return version, assets, nil
}
