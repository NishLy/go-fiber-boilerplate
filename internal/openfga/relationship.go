package fga

import (
	"github.com/NishLy/go-fiber-boilerplate/internal/platform/cache"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
)

func GetFGAWriteOptions(indentifier string) *client.ClientWriteOptions {
	modelID, _ := cache.Get[string](GetFGAClientModelKey(indentifier), 0, nil)
	var modelIDPtr *string

	if modelID != "" {
		modelIDPtr = openfga.PtrString(modelID)
	}

	var options = client.ClientWriteOptions{
		AuthorizationModelId: modelIDPtr,
	}
	return &options
}
