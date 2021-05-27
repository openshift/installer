package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

// WorkerPoolRequest provides worker pool data
// swagger:model
type WorkerPoolRequest struct {
	Cluster string `json:"cluster" description:"cluster name where the worker pool will be created"`
	WorkerPoolConfig
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
	Flavor      string            `json:"flavor"`
	ID          string            `json:"id"`
	Isolation   string            `json:"isolation"`
	Labels      map[string]string `json:"labels,omitempty"`
	Lifecycle   `json:"lifecycle"`
	VpcID       string     `json:"vpcID"`
	WorkerCount int        `json:"workerCount"`
	PoolName    string     `json:"poolName"`
	Provider    string     `json:"provider"`
	Zones       []ZoneResp `json:"zones"`
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

//Workers ...
type WorkerPool interface {
	CreateWorkerPool(workerPoolReq WorkerPoolRequest, target ClusterTargetHeader) (WorkerPoolResponse, error)
	GetWorkerPool(clusterNameOrID, workerPoolNameOrID string, target ClusterTargetHeader) (GetWorkerPoolResponse, error)
	ListWorkerPools(clusterNameOrID string, target ClusterTargetHeader) ([]GetWorkerPoolResponse, error)
	CreateWorkerPoolZone(workerPoolZone WorkerPoolZone, target ClusterTargetHeader) error
	DeleteWorkerPool(clusterNameOrID string, workerPoolNameOrID string, target ClusterTargetHeader) error
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
