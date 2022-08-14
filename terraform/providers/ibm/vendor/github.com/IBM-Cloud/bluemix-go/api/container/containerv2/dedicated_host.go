package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

// CreateDedicatedHostRequest provides dedicated host data for create call
// swagger:model
type CreateDedicatedHostRequest struct {
	Flavor     string `json:"flavor" description:""`
	HostPoolID string `json:"hostPoolID" description:""`
	Zone       string `json:"zone" description:""`
}

// CreateDedicatedHostResponse provides dedicated host id from create call
// swagger:model
type CreateDedicatedHostResponse struct {
	ID string `json:"id"`
}

// GetDedicatedHostResponse provides dedicated host data from get call
// swagger:model
type GetDedicatedHostResponse struct {
	Flavor           string                 `json:"flavor"`
	ID               string                 `json:"id"`
	Lifecycle        DedicatedHostLifecycle `json:"lifecycle"`
	PlacementEnabled bool                   `json:"placementEnabled"`
	Resources        DedicatedHostResources `json:"resources"`
	Workers          []DedicatedHostWorker  `json:"workers"`
	Zone             string                 `json:"zone"`
}

// DedicatedHostLifecycle ...
type DedicatedHostLifecycle struct {
	ActualState        string `json:"actualState"`
	DesiredState       string `json:"desiredState"`
	Message            string `json:"message"`
	MessageDate        string `json:"messageDate"`
	MessageDetails     string `json:"messageDetails"`
	MessageDetailsDate string `json:"messageDetailsDate"`
}

// DedicatedHostWorker ...
type DedicatedHostWorker struct {
	ClusterID    string `json:"clusterID"`
	Flavor       string `json:"flavor"`
	WorkerID     string `json:"workerID"`
	WorkerPoolID string `json:"workerPoolID"`
}

// DedicatedHostResources ...
type DedicatedHostResources struct {
	Capacity DedicatedHostResource `json:"capacity"`
	Consumed DedicatedHostResource `json:"consumed"`
}

// DedicatedHostResource ...
type DedicatedHostResource struct {
	MemoryBytes int64 `json:"memoryBytes"`
	VCPU        int64 `json:"vcpu"`
}

// RemoveDedicatedHostRequest provides dedicated host data for remove call
// swagger:model
type RemoveDedicatedHostRequest struct {
	HostID     string `json:"host" description:""`
	HostPoolID string `json:"hostPool" description:""`
}

// UpdateDedicatedHostPlacementRequest provides dedicated host data for update call
// swagger:model
type UpdateDedicatedHostPlacementRequest struct {
	HostPoolID string `json:"hostPoolID"`
	HostID     string `json:"hostID"`
}

//DedicatedHost ...
type DedicatedHost interface {
	CreateDedicatedHost(dedicatedHostReq CreateDedicatedHostRequest, target ClusterTargetHeader) (CreateDedicatedHostResponse, error)
	GetDedicatedHost(dedicatedHostID, dedicatedHostPoolID string, target ClusterTargetHeader) (GetDedicatedHostResponse, error)
	ListDedicatedHosts(dedicatedHostPoolID string, target ClusterTargetHeader) ([]GetDedicatedHostResponse, error)
	RemoveDedicatedHost(dedicatedHostReq RemoveDedicatedHostRequest, target ClusterTargetHeader) error
	EnableDedicatedHostPlacement(dedicatedHostReq UpdateDedicatedHostPlacementRequest, target ClusterTargetHeader) error
	DisableDedicatedHostPlacement(dedicatedHostReq UpdateDedicatedHostPlacementRequest, target ClusterTargetHeader) error
}

type dedicatedhost struct {
	client *client.Client
}

func newDedicatedHostAPI(c *client.Client) DedicatedHost {
	return &dedicatedhost{
		client: c,
	}
}

// GetDedicatedHost calls the API to list dedicated host s
func (w *dedicatedhost) ListDedicatedHosts(dedicatedHostPoolID string, target ClusterTargetHeader) ([]GetDedicatedHostResponse, error) {
	successV := []GetDedicatedHostResponse{}
	_, err := w.client.Get(fmt.Sprintf("/v2/getDedicatedHosts?dedicatedhostpool=%s", dedicatedHostPoolID), &successV, target.ToMap())
	return successV, err
}

// GetDedicatedHost calls the API to get a dedicated host
func (w *dedicatedhost) GetDedicatedHost(dedicatedHostID, dedicatedHostPoolID string, target ClusterTargetHeader) (GetDedicatedHostResponse, error) {
	var successV GetDedicatedHostResponse
	_, err := w.client.Get(fmt.Sprintf("/v2/getDedicatedHost?dedicatedhost=%s&dedicatedhostpool=%s", dedicatedHostID, dedicatedHostPoolID), &successV, target.ToMap())
	return successV, err
}

// CreateDedicatedHost calls the API to create a dedicated host
func (w *dedicatedhost) CreateDedicatedHost(createDedicatedHostReq CreateDedicatedHostRequest, target ClusterTargetHeader) (CreateDedicatedHostResponse, error) {
	var successV CreateDedicatedHostResponse
	_, err := w.client.Post("/v2/createDedicatedHost", createDedicatedHostReq, &successV, target.ToMap())
	return successV, err
}

// RemoveDedicatedHost calls the API to remove a dedicated host
func (w *dedicatedhost) RemoveDedicatedHost(removeDedicatedHostReq RemoveDedicatedHostRequest, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := w.client.Post("/v2/removeDedicatedHost", removeDedicatedHostReq, nil, target.ToMap())
	return err
}

// EnableDedicatedHostPlacement calls the API to enable placement on a dedicated host
func (w *dedicatedhost) EnableDedicatedHostPlacement(updateDedicatedHostReq UpdateDedicatedHostPlacementRequest, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := w.client.Post("/v2/enableDedicatedHostPlacement", updateDedicatedHostReq, nil, target.ToMap())
	return err
}

// DisableDedicatedHostPlacement calls the API to disable placement on a dedicated host
func (w *dedicatedhost) DisableDedicatedHostPlacement(updateDedicatedHostReq UpdateDedicatedHostPlacementRequest, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := w.client.Post("/v2/disableDedicatedHostPlacement", updateDedicatedHostReq, nil, target.ToMap())
	return err
}
