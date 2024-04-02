package grpcclient

type Config struct {
	AddressServiceUrl              string `toml:"address_service_url"`
	StorehouseServiceUrl           string `toml:"storehouse_service_url"`
	AutoReferenceCatalogServiceUrl string `toml:"auto_reference_catalog_service_url"`
}

func NewConfig() *Config {
	return &Config{
		AddressServiceUrl:              ":8081",
		StorehouseServiceUrl:           ":8082",
		AutoReferenceCatalogServiceUrl: ":8083",
	}
}
