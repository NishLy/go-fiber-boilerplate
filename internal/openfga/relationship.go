package fga

import (
	"context"

	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

func GetFGAWriteOptions(indentifier string) *client.ClientWriteOptions {
	var modelIDPtr *string
	var options = client.ClientWriteOptions{}

	if indentifier != "" {
		modelIDPtr = openfga.PtrString(indentifier)
		options.AuthorizationModelId = modelIDPtr
	}

	return &options
}

func WriteRelationships(clientIns *client.OpenFgaClient,
	options *client.ClientWriteOptions,
	relationships []client.ClientTupleKey) error {

	body := client.ClientWriteRequest{
		Writes: relationships,
	}

	_, err := clientIns.Write(context.Background()).Body(body).Options(*options).Execute()

	if err != nil {
		logger.Sugar.Errorf("Failed to write relationships: %v", err)
		return err
	}

	return err
}

func CheckRelationship(clientIns *client.OpenFgaClient, options *client.ClientCheckOptions, body client.ClientCheckRequest) (bool, error) {
	resp, err := clientIns.Check(context.Background()).Body(body).Options(*options).Execute()
	if err != nil {
		logger.Sugar.Errorf("Failed to check relationship: %v", err)
		return false, err
	}

	return *resp.Allowed, nil
}
