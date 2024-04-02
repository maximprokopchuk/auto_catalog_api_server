package config

import (
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/apiserver"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/grpcclient"
)

type Config struct {
	ApiServer  *apiserver.Config  `toml:"apiserver"`
	GrpcClient *grpcclient.Config `toml:"grpc_client"`
}

func NewConfig() *Config {
	return &Config{
		ApiServer:  apiserver.NewConfig(),
		GrpcClient: grpcclient.NewConfig(),
	}
}
