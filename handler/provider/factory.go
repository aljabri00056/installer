package provider

import (
	"fmt"
	"strings"
)

func NewProvider(providerType string, baseURL string) (Provider, error) {
	providerType = strings.ToLower(strings.TrimSpace(providerType))
	switch providerType {
	case "github", "":
		return &GitHub{BaseURL: "https://api.github.com"}, nil
	case "forgejo":
		if baseURL == "" {
			return nil, fmt.Errorf("baseURL is required for Forgejo provider")
		}
		return &GitHub{BaseURL: baseURL + "/api/v1"}, nil
	case "codeberg":
		return &GitHub{BaseURL: "https://codeberg.org/api/v1"}, nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
