package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

// CreateDedicatedHostPoolRequest provides dedicated host pool data for create call
// swagger:model
type CreateDedicatedHostPoolRequest struct {
	FlavorClass string `json:"flavorClass" description:""`
	Metro       string `json:"metro" description:""`
	Name        string `json:"name" description:""`
}

// CreateDedicatedHostPoolResponse provides dedicated host pool id from create call
// swagger:model
type CreateDedicatedHostPoolResponse struct {
	ID string `json:"id"`
}

// GetDedicatedHostPoolResponse provides dedicated host pool data from get call
// swagger:model
type GetDedicatedHostPoolResponse struct {
	FlavorClass string                        `json:"flavorClass"`
	HostCount   int64                         `json:"hostCount"`
	ID          string                        `json:"id"`
	Metro       string                        `json:"metro"`
	Name        string                        `json:"name"`
	State       string                        `json:"state"`
	WorkerPools []DedicatedHostPoolWorkerPool `json:"workerPools"`
	Zones       []DedicatedHostZoneResources  `json:"zones"`
}

// DedicatedHostPoolWorkerPool ...
type DedicatedHostPoolWorkerPool struct {
	ClusterID    string `json:"clusterID"`
	WorkerPoolID string `json:"workerPoolID"`
}

// DedicatedHostZoneResources ...
type DedicatedHostZoneResources struct {
	Capacity  DedicatedHostResource `json:"capacity"`
	HostCount int64                 `json:"hostCount"`
	Zone      string                `json:"zone"`
}

// RemoveDedicatedHostPoolRequest provides dedicated host pool data for remove call
// swagger:model
type RemoveDedicatedHostPoolRequest struct {
	HostPoolID string `json:"hostPool" description:""`
}

//DedicatedHostPool ...
type DedicatedHostPool interface {
	CreateDedicatedHostPool(dedicatedHostPoolReq CreateDedicatedHostPoolRequest, target ClusterTargetHeader) (CreateDedicatedHostPoolResponse, error)
	GetDedicatedHostPool(dedicatedHostPoolID string, target ClusterTargetHeader) (GetDedicatedHostPoolResponse, error)
	ListDedicatedHostPools(target ClusterTargetHeader) ([]GetDedicatedHostPoolResponse, error)
	RemoveDedicatedHostPool(dedicatedHostPoolReq RemoveDedicatedHostPoolRequest, target ClusterTargetHeader) error
}

type dedicatedhostpool struct {
	client *client.Client
}

func newDedicatedHostPoolAPI(c *client.Client) DedicatedHostPool {
	return &dedicatedhostpool{
		client: c,
	}
}

// GetDedicatedHostPool calls the API to list dedicated host pools
func (w *dedicatedhostpool) ListDedicatedHostPools(target ClusterTargetHeader) ([]GetDedicatedHostPoolResponse, error) {
	successV := []GetDedicatedHostPoolResponse{}
	_, err := w.client.Get("/v2/getDedicatedHostPools", &successV, target.ToMap())
	return successV, err
}

// GetDedicatedHostPool calls the API to get a dedicated host pool
func (w *dedicatedhostpool) GetDedicatedHostPool(dedicatedHostPoolID string, target ClusterTargetHeader) (GetDedicatedHostPoolResponse, error) {
	var successV GetDedicatedHostPoolResponse
	_, err := w.client.Get(fmt.Sprintf("/v2/getDedicatedHostPool?dedicatedhostpool=%s", dedicatedHostPoolID), &successV, target.ToMap())
	return successV, err
}

// CreateDedicatedHostPool calls the API to create a dedicated host pool
func (w *dedicatedhostpool) CreateDedicatedHostPool(createDedicatedHostPoolReq CreateDedicatedHostPoolRequest, target ClusterTargetHeader) (CreateDedicatedHostPoolResponse, error) {
	var successV CreateDedicatedHostPoolResponse
	_, err := w.client.Post("/v2/createDedicatedHostPool", createDedicatedHostPoolReq, &successV, target.ToMap())
	return successV, err
}

// RemoveDedicatedHostPool calls the API to remove a dedicated host pool
func (w *dedicatedhostpool) RemoveDedicatedHostPool(removeDedicatedHostPoolReq RemoveDedicatedHostPoolRequest, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := w.client.Post("/v2/removeDedicatedHostPool", removeDedicatedHostPoolReq, nil, target.ToMap())
	return err
}
