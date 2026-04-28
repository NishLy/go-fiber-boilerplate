package fga

import (
	"github.com/NishLy/go-fiber-boilerplate/config"
	"github.com/openfga/go-sdk/client"
)

func InitOpenFGA(indentifier string) (*client.OpenFgaClient, error) {
	cfg := config.Get()

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl:  cfg.OPEN_FGA_API_URL, // OpenFGA server address
		StoreId: indentifier,          // Created via CLI or API
	})

	return fgaClient, err
}
