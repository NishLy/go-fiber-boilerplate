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

	// change indentifier to ulid format: appname:tenantid

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl:  cfg.OPEN_FGA_API_URL, // OpenFGA server address
		StoreId: *indentifier,         // Created via CLI or API
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenFGA client: %w", err)
	}

	return fgaClient, nil
}

func CreateStore(ctx context.Context, identifier string) (*string, error) {
	cfg := config.Get()

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl: cfg.OPEN_FGA_API_URL, // OpenFGA server address
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenFGA client: %w", err)
	}

	store, err := fgaClient.CreateStore(ctx).
		Body(client.ClientCreateStoreRequest{
			Name: identifier,
		}).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create fga store: %w", err)
	}

	storeID := store.GetId()

	return &storeID, nil
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
	} else {
		name := fmt.Sprintf("%s:%s", config.Get().APP_NAME, identifier)
		storeID, err := CreateStore(context.Background(), name)
		if err != nil {
			return nil, fmt.Errorf("failed to provision new store: %w", err)
		}

		fgaClient, err = InitOpenFGA(storeID)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize OpenFGA client: %w", err)
		}
	}

	fgaClients[*storeID] = fgaClient
	return fgaClient, nil
}
