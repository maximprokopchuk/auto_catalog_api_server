package grpcclient

import (
	address_api "github.com/maximprokopchuk/address_service/pkg/api"
	auto_catalog_api "github.com/maximprokopchuk/auto_reference_catalog_service/pkg/api"
	storehouse_api "github.com/maximprokopchuk/storehouse_service/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	config            *Config
	AddressClient     address_api.AddressServiceClient
	StoreHouseClient  storehouse_api.StorehouseServiceClient
	AutoCatalogClient auto_catalog_api.AutoCatalogServiceClient
}

func New(config *Config) *GRPCClient {
	return &GRPCClient{
		config: config,
	}
}

func (client *GRPCClient) Init() error {
	address_client, err := client.initAddressClient()
	if err != nil {
		return nil
	}
	storehouse_client, err := client.initStorehouseClient()
	if err != nil {
		return nil
	}
	auto_catalog_client, err := client.initAutoCatalogServiceClient()
	if err != nil {
		return nil
	}
	client.AddressClient = address_client
	client.StoreHouseClient = storehouse_client
	client.AutoCatalogClient = auto_catalog_client
	return nil
}

func (client *GRPCClient) initAddressClient() (address_api.AddressServiceClient, error) {
	conn, err := grpc.Dial(client.config.AddressServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}
	address_client := address_api.NewAddressServiceClient(conn)
	return address_client, nil
}

func (client *GRPCClient) initStorehouseClient() (storehouse_api.StorehouseServiceClient, error) {
	conn, err := grpc.Dial(client.config.StorehouseServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}
	storehouse_client := storehouse_api.NewStorehouseServiceClient(conn)
	return storehouse_client, nil
}

func (client *GRPCClient) initAutoCatalogServiceClient() (auto_catalog_api.AutoCatalogServiceClient, error) {
	conn, err := grpc.Dial(client.config.AutoReferenceCatalogServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}
	auto_catalog_client := auto_catalog_api.NewAutoCatalogServiceClient(conn)
	return auto_catalog_client, nil
}
