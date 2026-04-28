package fga

import (
	"context"
	"fmt"

	"github.com/NishLy/go-fiber-boilerplate/config"
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	"github.com/oklog/ulid"
	"github.com/openfga/go-sdk/client"
)

var fgaClients = make(map[string]*client.OpenFgaClient)

func InitOpenFGA(indentifier *string) (*client.OpenFgaClient, error) {
	cfg := config.Get()

	fmt.Printf("Initializing OpenFGA client for store ID: %s\n", *indentifier)

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl:  cfg.OPEN_FGA_API_URL, // OpenFGA server address
		StoreId: *indentifier,         // Created via CLI or API
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenFGA client: %w", err)
	}

	return fgaClient, nil
}

func isValidULID(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

func GetFGAClient(identifier string) (*client.OpenFgaClient, error) {
	if client, exists := fgaClients[identifier]; exists {
		return client, nil
	}

	logger.Sugar.Infof("Initialized new OpenFGA client for identifier %s", identifier)
	var fgaClient *client.OpenFgaClient
	var storeID *string

	// If the identifier is a valid ULID, use it directly as the store ID. Otherwise, treat it as a tenant identifier and generate a store ID based on the app name and tenant ID.
	if isValidULID(identifier) {
		storeID = &identifier
		fgaClient, err := InitOpenFGA(storeID)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize OpenFGA client: %w", err)
		}
		fgaClients[*storeID] = fgaClient
		return fgaClient, nil
	} else {
		storeID, err := CreateStore(context.Background(), identifier)
		if err != nil {
			return nil, fmt.Errorf("failed to provision new store: %w", err)
		}

		fmt.Printf("Provisioned new store with ID: %s\n", *storeID)
		fgaClient, err = InitOpenFGA(storeID)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize OpenFGA client: %w", err)
		}
		fmt.Printf("FGA %v client initialized successfully for store ID: %s\n", fgaClient, *storeID)
		fgaClients[*storeID] = fgaClient
		return fgaClient, nil
	}
}
