package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RepoInfo struct {
	Private bool
}

type Asset struct {
	Size        int
	Name        string
	OS          string
	Arch        string
	Type        string
	URL         string
	DownloadURL string
}

func (a Asset) Key() string {
	return a.OS + "/" + a.Arch
}

func (a Asset) DisplayKey() string {
	os := a.OS
	if os == "darwin" {
		os = "macOS"
	}
	return os + "/" + a.Arch
}

func (a Asset) Is32Bit() bool {
	return a.Arch == "386"
}

func (a Asset) IsMac() bool {
	return a.OS == "darwin"
}

func (a Asset) IsMacM1() bool {
	return a.IsMac() && a.Arch == "arm64"
}

type Provider interface {
	GetRepo(user, repo, token string) (*RepoInfo, error)
	GetReleaseAssets(user, repo, release, token string) (string, []Asset, error)
}

type BaseProvider struct{}

func (p *BaseProvider) get(url string, token string, v any) error {
	req, _ := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %s: %s", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return fmt.Errorf("not found: url %s", url)
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s %s", http.StatusText(resp.StatusCode), string(b))
	}
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("decode failed: %s: %s", url, err)
		}
	}
	return nil
}
