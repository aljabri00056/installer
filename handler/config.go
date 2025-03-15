package handler

type Config struct {
	Host        string `opts:"help=host, env=HTTP_HOST"`
	Port        int    `opts:"help=port, env"`
	User        string `opts:"help=default user when not provided in URL, env"`
	Provider    string `opts:"help=git provider (github,codeberg,forgejo), env=GIT_PROVIDER"`
	ProviderURL string `opts:"help=base URL for forgejo/gitea instance, env=PROVIDER_URL"`
}

var DefaultConfig = Config{
	Port: 8080,
}
