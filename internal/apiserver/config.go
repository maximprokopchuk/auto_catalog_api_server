package apiserver

type Config struct {
	BindAddr      string `toml:"bind_url"`
	AllowedOrigin string `toml:"allowed_origin"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8081",
	}
}
