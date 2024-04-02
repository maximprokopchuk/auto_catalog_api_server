package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/apiserver"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/config"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config_path", "configs/config.toml", "path to config TOML file")
}

func run() error {
	cfg := config.NewConfig()
	_, err := toml.DecodeFile(configPath, cfg)

	if err != nil {
		return err
	}

	fmt.Println(cfg.ApiServer.BindAddr)

	s := apiserver.New(cfg.ApiServer)

	if err := s.Start(); err != nil {
		return nil
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
