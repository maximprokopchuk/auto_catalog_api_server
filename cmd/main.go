package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/apiserver"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/config"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/grpcclient"
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

	grpc_client := grpcclient.New(cfg.GrpcClient)
	err = grpc_client.Init()

	if err != nil {
		return err
	}

	s := apiserver.New(cfg.ApiServer, grpc_client)

	log.Println("LISTEN " + cfg.ApiServer.BindAddr)

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
