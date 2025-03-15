package provider

import (
	"fmt"
)

type GitLab struct {
	BaseProvider
	BaseURL string
}

type glAsset struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Size int    `json:"size"`
}

type glRelease struct {
	Name        string    `json:"name"`
	TagName     string    `json:"tag_name"`
	Assets      struct {
		Links []glAsset `json:"links"`
	} `json:"assets"`
}

type glRepo struct {
	Private bool `json:"visibility"`
}

func (g *GitLab) GetRepo(user, repo, token string) (*RepoInfo, error) {
	if user == "" || repo == "" {
		return nil, fmt.Errorf("user and repo are required")
	}

	url := fmt.Sprintf("%s/projects/%s%%2F%s", g.BaseURL, user, repo)
	var res glRepo
	if err := g.get(url, token, &res); err != nil {
		return nil, fmt.Errorf("failed to get repo info: %w", err)
	}
	return &RepoInfo{Private: res.Private}, nil
}

func (g *GitLab) GetReleaseAssets(user, repo, release, token string) (string, []Asset, error) {
	var assets []Asset
	var version string

	url := fmt.Sprintf("%s/projects/%s%%2F%s/releases", g.BaseURL, user, repo)
	
	if release == "" || release == "latest" {
		var releases []glRelease
		if err := g.get(url, token, &releases); err != nil {
			return "", nil, err
		}
		if len(releases) == 0 {
			return "", nil, fmt.Errorf("no releases found")
		}
		
		version = releases[0].TagName
		for _, a := range releases[0].Assets.Links {
			assets = append(assets, Asset{
				Name:        a.Name,
				URL:         a.URL,
				Size:        a.Size,
				DownloadURL: a.URL,
			})
		}
	} else {
		version = release
		url = fmt.Sprintf("%s/projects/%s%%2F%s/releases/%s", g.BaseURL, user, repo, release)
		var resp glRelease
		if err := g.get(url, token, &resp); err != nil {
			return "", nil, err
		}
		
		for _, a := range resp.Assets.Links {
			assets = append(assets, Asset{
				Name:        a.Name,
				URL:         a.URL,
				Size:        a.Size,
				DownloadURL: a.URL,
			})
		}
	}

	return version, assets, nil
}
