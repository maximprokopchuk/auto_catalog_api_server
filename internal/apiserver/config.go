package apiserver

type Config struct {
	BindAddr string `toml:"bind_url"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8081",
	}
}
