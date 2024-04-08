package grpcclient_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/config"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/grpcclient"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	cfg := config.NewConfig()
	_, err := toml.DecodeFile("../../configs/config.test.toml", cfg)
	assert.Nil(t, err)
	client := grpcclient.New(cfg.GrpcClient)
	err = client.Init()
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.AddressClient)
	assert.NotNil(t, client.StoreHouseClient)
	assert.NotNil(t, client.AutoReferenceCatalogClient)
}
