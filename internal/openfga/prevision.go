package fga

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/NishLy/go-fiber-boilerplate/config"
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/language/pkg/go/transformer"
)

func CreateStore(ctx context.Context, identifier string) (*string, *string, error) {
	cfg := config.Get()

	logger.Sugar.Infof("Creating new OpenFGA store with identifier: %s", identifier)
	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl: cfg.OPEN_FGA_API_URL, // OpenFGA server address
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize OpenFGA client: %w", err)
	}

	store, err := fgaClient.CreateStore(ctx).
		Body(client.ClientCreateStoreRequest{
			Name: identifier,
		}).Execute()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create fga store: %w", err)
	}

	storeID := store.GetId()
	// init models for the store
	modelID, err := MigrateFromFolder(ctx, storeID, cfg.OPEN_FGA_MODEL_DIR)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to migrate models: %w", err)
	}

	return &storeID, &modelID, nil
}

func MigrateFromFolder(ctx context.Context, storeID string, folderPath string) (string, error) {

	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl:  config.Get().OPEN_FGA_API_URL, // OpenFGA server address
		StoreId: storeID,                       // Created via CLI or API
	})

	if err != nil {
		return "", fmt.Errorf("failed to initialize OpenFGA client: %w", err)
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return "", fmt.Errorf("failed to read openfga folder %q: %w", folderPath, err)
	}

	var dslParts []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".fga") {
			continue
		}

		filePath := filepath.Join(folderPath, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read file %q: %w", filePath, err)
		}
		dslParts = append(dslParts, string(data))
	}

	if len(dslParts) == 0 {
		return "", fmt.Errorf("no .fga files found in %q", folderPath)
	}

	// 2. Join all DSL parts and transform to an authorization model
	combinedDSL := strings.Join(dslParts, "\n\n")

	authModel, err := transformer.TransformDSLToProto(combinedDSL)
	if err != nil {
		return "", fmt.Errorf("failed to parse DSL: %w", err)
	}

	// 3. Convert proto TypeDefinitions to SDK TypeDefinitions
	typeDefinitions := make([]openfga.TypeDefinition, 0, len(authModel.GetTypeDefinitions()))
	for _, td := range authModel.GetTypeDefinitions() {
		typeDefinitions = append(typeDefinitions, *openfga.NewTypeDefinition(td.GetType()))
	}

	openFGAConditions := make(map[string]openfga.Condition)
	for name, condition := range authModel.GetConditions() {
		openFGAConditions[name] = openfga.Condition{
			Expression: condition.GetExpression(),
		}
	}

	body := client.ClientWriteAuthorizationModelRequest{
		SchemaVersion:   authModel.GetSchemaVersion(),
		TypeDefinitions: typeDefinitions,
		Conditions:      &openFGAConditions,
	}

	// 4. Write the merged model to the store
	resp, err := fgaClient.WriteAuthorizationModel(ctx).Body(body).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to write authorization model: %w", err)
	}

	modelID := resp.GetAuthorizationModelId()
	return modelID, nil
}
