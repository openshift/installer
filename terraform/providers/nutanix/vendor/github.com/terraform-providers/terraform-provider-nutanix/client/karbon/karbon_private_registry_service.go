package karbon

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// PrivateRegistryOperations ...
type PrivateRegistryOperations struct {
	client *client.Client
}

// Service ...
type PrivateRegistryService interface {
	// karbon v2.1
	ListKarbonPrivateRegistries() (*PrivateRegistryListResponse, error)
	CreateKarbonPrivateRegistry(createRequest *PrivateRegistryIntentInput) (*PrivateRegistryResponse, error)
	GetKarbonPrivateRegistry(name string) (*PrivateRegistryResponse, error)
	DeleteKarbonPrivateRegistry(name string) (*PrivateRegistryOperationResponse, error)
}

func (op PrivateRegistryOperations) ListKarbonPrivateRegistries() (*PrivateRegistryListResponse, error) {
	ctx := context.TODO()
	path := "/v1-alpha.1/registries"
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonPrivateRegistryListResponse := new(PrivateRegistryListResponse)
	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryListResponse, op.client.Do(ctx, req, karbonPrivateRegistryListResponse)
}

func (op PrivateRegistryOperations) CreateKarbonPrivateRegistry(createRequest *PrivateRegistryIntentInput) (*PrivateRegistryResponse, error) {
	ctx := context.TODO()
	path := "/v1-alpha.1/registries"
	req, err := op.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	karbonPrivateRegistryResponse := new(PrivateRegistryResponse)
	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryResponse, op.client.Do(ctx, req, karbonPrivateRegistryResponse)
}

func (op PrivateRegistryOperations) GetKarbonPrivateRegistry(name string) (*PrivateRegistryResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1-alpha.1/registries/%s", name)
	fmt.Printf("Path: %s", path)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonPrivateRegistryResponse := new(PrivateRegistryResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryResponse, op.client.Do(ctx, req, karbonPrivateRegistryResponse)
}

func (op PrivateRegistryOperations) DeleteKarbonPrivateRegistry(name string) (*PrivateRegistryOperationResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/v1-alpha.1/registries/%s", name)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	karbonPrivateRegistryOperationResponse := new(PrivateRegistryOperationResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryOperationResponse, op.client.Do(ctx, req, karbonPrivateRegistryOperationResponse)
}
