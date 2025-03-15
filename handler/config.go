package handler

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port            int               `opts:"help=port, env=HTTP_PORT"`
	User            string            `opts:"help=default user when not provided in URL, env=DEFAULT_USER"`
	Provider        string            `opts:"help=git provider (github,codeberg,forgejo), env=GIT_PROVIDER"`
	ProviderURL     string            `opts:"help=base URL for forgejo/gitea instance, env=PROVIDER_URL"`
	LogLevel        string            `opts:"help=log level (debug,info,warn,error), env=LOG_LEVEL"`
	RepoProviderMap map[string]string `opts:"help=repository to provider mapping"`
}

var DefaultConfig = Config{
	Port:     8080,
	LogLevel: "info",
}

func GetConfigFromEnv() Config {
	config := DefaultConfig

	if port := getEnv("HTTP_PORT", "8080"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Port = p
		}
	}
	if user := getEnv("DEFAULT_USER", ""); user != "" {
		config.User = user
	}
	if provider := getEnv("GIT_PROVIDER", ""); provider != "" {
		config.Provider = provider
	}
	if providerURL := getEnv("PROVIDER_URL", ""); providerURL != "" {
		config.ProviderURL = providerURL
	}
	if logLevel := getEnv("LOG_LEVEL", "info"); logLevel != "" {
		config.LogLevel = logLevel
	}

	config.RepoProviderMap = make(map[string]string)
	if mapStr := getEnv("REPO_PROVIDER_MAP", ""); mapStr != "" {
		for _, mapping := range strings.Split(mapStr, ",") {
			parts := strings.Split(mapping, "=")
			if len(parts) == 2 {
				config.RepoProviderMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	return config
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
