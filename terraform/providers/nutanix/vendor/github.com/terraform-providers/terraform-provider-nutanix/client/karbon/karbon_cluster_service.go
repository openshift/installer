package karbon

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// ClusterOperations ...
type ClusterOperations struct {
	client *client.Client
}

// Service ...
type ClusterService interface {
	// karbon v2.1
	ListKarbonClusters() (*ClusterListIntentResponse, error)
	CreateKarbonCluster(createRequest *ClusterIntentInput) (*ClusterActionResponse, error)
	GetKarbonCluster(karbonClusterName string) (*ClusterIntentResponse, error)
	GetKarbonClusterNodePool(karbonClusterName string, nodePoolName string) (*ClusterNodePool, error)
	DeleteKarbonCluster(karbonClusterName string) (*ClusterActionResponse, error)
	GetKubeConfigForKarbonCluster(karbonClusterName string) (*ClusterKubeconfigResponse, error)
	GetSSHConfigForKarbonCluster(karbonClusterName string) (*ClusterSSHconfig, error)
	ScaleUpKarbonCluster(karbonClusterName, karbonNodepoolName string, scaleUpRequest *ClusterScaleUpIntentInput) (*ClusterActionResponse, error)
	ScaleDownKarbonCluster(karbonClusterName, karbonNodepoolName string, scaleDownRequest *ClusterScaleDownIntentInput) (*ClusterActionResponse, error)
	// registries
	ListPrivateRegistries(karbonClusterName string) (*PrivateRegistryListResponse, error)
	AddPrivateRegistry(karbonClusterName string, createRequest PrivateRegistryOperationIntentInput) (*PrivateRegistryResponse, error)
	DeletePrivateRegistry(karbonClusterName string, privateRegistryName string) (*PrivateRegistryOperationResponse, error)
}

// karbon 2.1
func (op ClusterOperations) ListKarbonClusters() (*ClusterListIntentResponse, error) {
	ctx := context.TODO()
	path := "/v1-beta.1/k8s/clusters"
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterListIntentResponse := new(ClusterListIntentResponse)
	if err != nil {
		return nil, err
	}

	return karbonClusterListIntentResponse, op.client.Do(ctx, req, karbonClusterListIntentResponse)
}

func (op ClusterOperations) CreateKarbonCluster(createRequest *ClusterIntentInput) (*ClusterActionResponse, error) {
	ctx := context.TODO()

	path := "/v1/k8s/clusters"
	req, err := op.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	karbonClusterActionResponse := new(ClusterActionResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterActionResponse, op.client.Do(ctx, req, karbonClusterActionResponse)
}

func (op ClusterOperations) GetKarbonCluster(name string) (*ClusterIntentResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1/k8s/clusters/%s", name)
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterIntentResponse := new(ClusterIntentResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterIntentResponse, op.client.Do(ctx, req, karbonClusterIntentResponse)
}

func (op ClusterOperations) GetKarbonClusterNodePool(name string, nodePoolName string) (*ClusterNodePool, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1-beta.1/k8s/clusters/%s/node-pools/%s", name, nodePoolName)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterNodePool := new(ClusterNodePool)

	if err != nil {
		return nil, err
	}

	return karbonClusterNodePool, op.client.Do(ctx, req, karbonClusterNodePool)
}

func (op ClusterOperations) DeleteKarbonCluster(name string) (*ClusterActionResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1/k8s/clusters/%s", name)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	karbonClusterActionResponse := new(ClusterActionResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterActionResponse, op.client.Do(ctx, req, karbonClusterActionResponse)
}

func (op ClusterOperations) GetKubeConfigForKarbonCluster(name string) (*ClusterKubeconfigResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1/k8s/clusters/%s/kubeconfig", name)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterKubeconfigResponse := new(ClusterKubeconfigResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterKubeconfigResponse, op.client.Do(ctx, req, karbonClusterKubeconfigResponse)
}

func (op ClusterOperations) GetSSHConfigForKarbonCluster(name string) (*ClusterSSHconfig, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1/k8s/clusters/%s/ssh", name)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonClusterSSHconfig := new(ClusterSSHconfig)

	if err != nil {
		return nil, err
	}

	return karbonClusterSSHconfig, op.client.Do(ctx, req, karbonClusterSSHconfig)
}

func (op ClusterOperations) ListPrivateRegistries(karbonClusterName string) (*PrivateRegistryListResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/v1-alpha.1/k8s/clusters/%s/registries", karbonClusterName)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonPrivateRegistryListResponse := new(PrivateRegistryListResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryListResponse, op.client.Do(ctx, req, karbonPrivateRegistryListResponse)
}

func (op ClusterOperations) AddPrivateRegistry(karbonClusterName string, createRequest PrivateRegistryOperationIntentInput) (*PrivateRegistryResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/v1-alpha.1/k8s/clusters/%s/registries", karbonClusterName)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	karbonPrivateRegistryResponse := new(PrivateRegistryResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryResponse, op.client.Do(ctx, req, karbonPrivateRegistryResponse)
}

func (op ClusterOperations) DeletePrivateRegistry(karbonClusterName string, privateRegistryName string) (*PrivateRegistryOperationResponse, error) {
	ctx := context.TODO()
	path := fmt.Sprintf("/v1-alpha.1/k8s/clusters/%s/registries/%s", karbonClusterName, privateRegistryName)

	req, err := op.client.NewRequest(ctx, http.MethodDelete, path, nil)
	karbonPrivateRegistryOperationResponse := new(PrivateRegistryOperationResponse)

	if err != nil {
		return nil, err
	}

	return karbonPrivateRegistryOperationResponse, op.client.Do(ctx, req, karbonPrivateRegistryOperationResponse)
}

func (op ClusterOperations) ScaleUpKarbonCluster(karbonClusterName, karbonNodepoolName string, scaleUpRequest *ClusterScaleUpIntentInput) (*ClusterActionResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1-alpha.1/k8s/clusters/%s/node-pools/%s/add-nodes", karbonClusterName, karbonNodepoolName)
	req, err := op.client.NewRequest(ctx, http.MethodPost, path, scaleUpRequest)
	karbonClusterActionResponse := new(ClusterActionResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterActionResponse, op.client.Do(ctx, req, karbonClusterActionResponse)
}

func (op ClusterOperations) ScaleDownKarbonCluster(karbonClusterName, karbonNodepoolName string, scaleDownRequest *ClusterScaleDownIntentInput) (*ClusterActionResponse, error) {
	ctx := context.TODO()

	path := fmt.Sprintf("/v1-alpha.1/k8s/clusters/%s/node-pools/%s/remove-nodes", karbonClusterName, karbonNodepoolName)
	req, err := op.client.NewRequest(ctx, http.MethodPost, path, scaleDownRequest)
	karbonClusterActionResponse := new(ClusterActionResponse)

	if err != nil {
		return nil, err
	}

	return karbonClusterActionResponse, op.client.Do(ctx, req, karbonClusterActionResponse)
}
