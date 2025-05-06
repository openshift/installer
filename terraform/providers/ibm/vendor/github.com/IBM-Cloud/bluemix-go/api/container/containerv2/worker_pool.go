package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

// CommonWorkerPoolConfig provides common worker pool data for cluster and workerpool operations
type CommonWorkerPoolConfig struct {
	DiskEncryption         bool                    `json:"diskEncryption,omitempty"`
	Entitlement            string                  `json:"entitlement"`
	Flavor                 string                  `json:"flavor"`
	Isolation              string                  `json:"isolation,omitempty"`
	Labels                 map[string]string       `json:"labels,omitempty"`
	Name                   string                  `json:"name" binding:"required" description:"The workerpool's name"`
	OperatingSystem        string                  `json:"operatingSystem,omitempty"`
	VpcID                  string                  `json:"vpcID"`
	WorkerCount            int                     `json:"workerCount"`
	Zones                  []Zone                  `json:"zones"`
	WorkerVolumeEncryption *WorkerVolumeEncryption `json:"workerVolumeEncryption,omitempty"`
	SecondaryStorageOption string                  `json:"secondaryStorageOption,omitempty"`
}

// WorkerPoolRequest provides worker pool data
// swagger:model
type WorkerPoolRequest struct {
	Cluster    string `json:"cluster" description:"cluster name where the worker pool will be created"`
	HostPoolID string `json:"hostPool,omitempty"`
	CommonWorkerPoolConfig
}
type WorkerPoolTaintRequest struct {
	Cluster    string            `json:"cluster" description:"cluster name"`
	WorkerPool string            `json:"workerpool" description:"worker Pool name"`
	Taints     map[string]string `json:"taints" description:"map of taints that has to be applied on workerpool"`
}

// WorkerPoolResponse provides worker pool data
// swagger:model
type WorkerPoolResponse struct {
	ID string `json:"workerPoolID"`
}

type WorkerPoolZone struct {
	Cluster      string `json:"cluster"`
	Id           string `json:"id"`
	SubnetID     string `json:"subnetID"`
	WorkerPoolID string `json:"workerPoolID"`
}

type GetWorkerPoolResponse struct {
	AutoscaleEnabled       bool              `json:"autoscaleEnabled,omitempty"`
	HostPoolID             string            `json:"dedicatedHostPoolId,omitempty"`
	Flavor                 string            `json:"flavor"`
	ID                     string            `json:"id"`
	Isolation              string            `json:"isolation"`
	Labels                 map[string]string `json:"labels,omitempty"`
	Lifecycle              `json:"lifecycle"`
	OperatingSystem        string                  `json:"operatingSystem,omitempty"`
	PoolName               string                  `json:"poolName"`
	Provider               string                  `json:"provider"`
	SecondaryStorageOption *DiskConfigResp         `json:"secondaryStorageOption,omitempty"`
	Taints                 map[string]string       `json:"taints,omitempty"`
	VpcID                  string                  `json:"vpcID"`
	WorkerCount            int                     `json:"workerCount"`
	WorkerVolumeEncryption *WorkerVolumeEncryption `json:"workerVolumeEncryption,omitempty"`
	Zones                  []ZoneResp              `json:"zones"`
}

// DiskConfigResp response type for describing a disk configuration
// swagger:model
type DiskConfigResp struct {
	Name  string `json:"name,omitempty"`
	Count int
	// the size of each individual device in GB
	Size              int
	DeviceType        string
	RAIDConfiguration string
	Profile           string `json:"profile,omitempty"`
}

type Lifecycle struct {
	ActualState  string `json:"actualState"`
	DesiredState string `json:"desiredState"`
}

type ZoneResp struct {
	ID          string   `json:"id"`
	WorkerCount int      `json:"workerCount"`
	Subnets     []Subnet `json:"subnets"`
}

type Subnet struct {
	ID      string `json:"id"`
	Primary bool   `json:"primary"`
}

type ResizeWorkerPoolReq struct {
	Cluster    string `json:"cluster"`
	Size       int64  `json:"size"`
	Workerpool string `json:"workerpool"`
}

//Workers ...
type WorkerPool interface {
	CreateWorkerPool(workerPoolReq WorkerPoolRequest, target ClusterTargetHeader) (WorkerPoolResponse, error)
	GetWorkerPool(clusterNameOrID, workerPoolNameOrID string, target ClusterTargetHeader) (GetWorkerPoolResponse, error)
	ListWorkerPools(clusterNameOrID string, target ClusterTargetHeader) ([]GetWorkerPoolResponse, error)
	CreateWorkerPoolZone(workerPoolZone WorkerPoolZone, target ClusterTargetHeader) error
	DeleteWorkerPool(clusterNameOrID string, workerPoolNameOrID string, target ClusterTargetHeader) error
	UpdateWorkerPoolTaints(taintRequest WorkerPoolTaintRequest, target ClusterTargetHeader) error
	ResizeWorkerPool(resizeWorkerPoolReq ResizeWorkerPoolReq, target ClusterTargetHeader) error
}

type workerpool struct {
	client *client.Client
}

func newWorkerPoolAPI(c *client.Client) WorkerPool {
	return &workerpool{
		client: c,
	}
}

// GetWorkerPool calls the API to get a worker pool
func (w *workerpool) ListWorkerPools(clusterNameOrID string, target ClusterTargetHeader) ([]GetWorkerPoolResponse, error) {
	successV := []GetWorkerPoolResponse{}
	_, err := w.client.Get(fmt.Sprintf("/v2/vpc/getWorkerPools?cluster=%s", clusterNameOrID), &successV, target.ToMap())
	return successV, err
}

// GetWorkerPool calls the API to get a worker pool
func (w *workerpool) GetWorkerPool(clusterNameOrID, workerPoolNameOrID string, target ClusterTargetHeader) (GetWorkerPoolResponse, error) {
	var successV GetWorkerPoolResponse
	_, err := w.client.Get(fmt.Sprintf("/v2/vpc/getWorkerPool?cluster=%s&workerpool=%s", clusterNameOrID, workerPoolNameOrID), &successV, target.ToMap())
	return successV, err
}

// CreateWorkerPool calls the API to create a worker pool
func (w *workerpool) CreateWorkerPool(workerPoolReq WorkerPoolRequest, target ClusterTargetHeader) (WorkerPoolResponse, error) {
	var successV WorkerPoolResponse
	_, err := w.client.Post("/v2/vpc/createWorkerPool", workerPoolReq, &successV, target.ToMap())
	return successV, err
}

// DeleteWorkerPool calls the API to remove a worker pool
func (w *workerpool) DeleteWorkerPool(clusterNameOrID string, workerPoolNameOrID string, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := w.client.Delete(fmt.Sprintf("/v1/clusters/%s/workerpools/%s", clusterNameOrID, workerPoolNameOrID), target.ToMap())
	return err
}

// CreateWorkerPoolZone calls the API to add a zone to a cluster and worker pool
func (w *workerpool) CreateWorkerPoolZone(workerPoolZone WorkerPoolZone, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := w.client.Post("/v2/vpc/createWorkerPoolZone", workerPoolZone, nil, target.ToMap())
	return err
}

// UpdateWorkerPoolTaints calls the API to update taints to a worker pool
func (w *workerpool) UpdateWorkerPoolTaints(taintRequest WorkerPoolTaintRequest, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := w.client.Post("/v2/setWorkerPoolTaints", taintRequest, nil, target.ToMap())
	return err
}

// ResizeWorkerPool calls the API to resize an existing worker pool.
func (w *workerpool) ResizeWorkerPool(resizeWorkerPoolReq ResizeWorkerPoolReq, target ClusterTargetHeader) error {
	// Make the request, don't care about return value
	_, err := w.client.Post("/v2/resizeWorkerPool", resizeWorkerPoolReq, nil, target.ToMap())
	return err
}
