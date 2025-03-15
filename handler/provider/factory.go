package provider

import (
	"fmt"
	"strings"
)

const (
	DefaultGitHubAPI = "https://api.github.com"
	DefaultCodebergAPI = "https://codeberg.org/api/v1"
)

// NewProvider creates a new provider instance based on the provider type
func NewProvider(providerType string, baseURL string) (Provider, error) {
	providerType = strings.ToLower(strings.TrimSpace(providerType))
	
	switch providerType {
	case "github", "":
		return &GitHub{BaseURL: DefaultGitHubAPI}, nil
	case "forgejo":
		if baseURL == "" {
			return nil, fmt.Errorf("baseURL is required for Forgejo provider")
		}
		return &GitHub{BaseURL: fmt.Sprintf("%s/api/v1", strings.TrimSuffix(baseURL, "/"))}, nil
	case "codeberg": 
		return &GitHub{BaseURL: DefaultCodebergAPI}, nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s (supported: github, codeberg, forgejo)", providerType)
	}
}
