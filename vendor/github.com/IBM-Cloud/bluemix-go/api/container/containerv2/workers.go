package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

//Worker ...
type Worker struct {
	Billing           string `json:"billing,omitempty"`
	Flavor            string `json:"flavor"`
	ID                string `json:"id"`
	KubeVersion       KubeDetails
	Location          string          `json:"location"`
	PoolID            string          `json:"poolid"`
	PoolName          string          `json:"poolName"`
	LifeCycle         WorkerLifeCycle `json:"lifecycle"`
	Health            HealthStatus    `json:"health"`
	NetworkInterfaces []Network       `json:"networkInterfaces"`
}

type KubeDetails struct {
	Actual    string `json:"actual"`
	Desired   string `json:"desired"`
	Eos       string `json:"eos"`
	MasterEOS string `json:"masterEos"`
	Target    string `json:"target"`
}
type HealthStatus struct {
	Message string `json:"message"`
	State   string `json:"state"`
}
type WorkerLifeCycle struct {
	ReasonForDelete    string `json:"reasonForDelete"`
	ActualState        string `json:"actualState"`
	DesiredState       string `json:"desiredState"`
	Message            string `json:"message"`
	MessageDate        string `json:"messageDate"`
	MessageDetails     string `json:"messageDetails"`
	MessageDetailsDate string `json:"messageDetailsDate"`
	PendingOperation   string `json:"pendingOperation"`
}

type Network struct {
	Cidr      string `json:"cidr"`
	IpAddress string `json:"ipAddress"`
	Primary   bool   `json:"primary"`
	SubnetID  string `json:"subnetID"`
}

type ReplaceWorker struct {
	ClusterIDOrName string `json:"cluster"`
	Update          bool   `json:"update"`
	WorkerID        string `json:"workerID"`
}

type VoulemeAttachments struct {
	VolumeAttachments []VoulemeAttachment `json:"volume_attachments"`
}

type VoulemeAttachment struct {
	Id     string     `json:"id"`
	Volume Volume     `json:"volume"`
	Device DeviceInfo `json:"device"`
	Name   string     `json:"name"`
	Status string     `json:"status"`
	Type   string     `json:"type"`
}

type Volume struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type DeviceInfo struct {
	Id string `json:"id"`
}

type VolumeRequest struct {
	Cluster            string `json:"cluster"`
	VolumeAttachmentID string `json:"volumeAttachmentID,omitempty"`
	VolumeID           string `json:"volumeID,omitempty"`
	Worker             string `json:"worker"`
}

//Workers ...
type Workers interface {
	ListByWorkerPool(clusterIDOrName, workerPoolIDOrName string, showDeleted bool, target ClusterTargetHeader) ([]Worker, error)
	ListWorkers(clusterIDOrName string, showDeleted bool, target ClusterTargetHeader) ([]Worker, error)
	Get(clusterIDOrName, workerID string, target ClusterTargetHeader) (Worker, error)
	ReplaceWokerNode(clusterIDOrName, workerID string, target ClusterTargetHeader) (string, error)
	ListStorageAttachemnts(clusterIDOrName, workerID string, target ClusterTargetHeader) (VoulemeAttachments, error)
	GetStorageAttachment(clusterIDOrName, workerID, volumeAttachmentID string, target ClusterTargetHeader) (VoulemeAttachment, error)
	CreateStorageAttachment(payload VolumeRequest, target ClusterTargetHeader) (VoulemeAttachment, error)
	DeleteStorageAttachment(payload VolumeRequest, target ClusterTargetHeader) (string, error)
}

type worker struct {
	client *client.Client
}

func newWorkerAPI(c *client.Client) Workers {
	return &worker{
		client: c,
	}
}

//ListByWorkerPool ...
func (r *worker) ListByWorkerPool(clusterIDOrName, workerPoolIDOrName string, showDeleted bool, target ClusterTargetHeader) ([]Worker, error) {
	rawURL := fmt.Sprintf("/v2/vpc/getWorkers?cluster=%s&showDeleted=%t", clusterIDOrName, showDeleted)
	if len(workerPoolIDOrName) > 0 {
		rawURL += "&pool=" + workerPoolIDOrName
	}
	workers := []Worker{}
	_, err := r.client.Get(rawURL, &workers, target.ToMap())
	if err != nil {
		return nil, err
	}
	return workers, err
}

//ListWorkers ...
func (r *worker) ListWorkers(clusterIDOrName string, showDeleted bool, target ClusterTargetHeader) ([]Worker, error) {
	rawURL := fmt.Sprintf("/v2/vpc/getWorkers?cluster=%s&showDeleted=%t", clusterIDOrName, showDeleted)
	workers := []Worker{}
	_, err := r.client.Get(rawURL, &workers, target.ToMap())
	if err != nil {
		return nil, err
	}
	return workers, err
}

//Get ...
func (r *worker) Get(clusterIDOrName, workerID string, target ClusterTargetHeader) (Worker, error) {
	rawURL := fmt.Sprintf("/v2/vpc/getWorker?cluster=%s&worker=%s", clusterIDOrName, workerID)
	worker := Worker{}
	_, err := r.client.Get(rawURL, &worker, target.ToMap())
	if err != nil {
		return worker, err
	}
	return worker, err
}

func (r *worker) ReplaceWokerNode(clusterIDOrName, workerID string, target ClusterTargetHeader) (string, error) {
	payload := ReplaceWorker{
		ClusterIDOrName: clusterIDOrName,
		WorkerID:        workerID,
		Update:          true,
	}
	var response string
	_, err := r.client.Post("/v2/vpc/replaceWorker", payload, &response, target.ToMap())
	if err != nil {
		return response, err
	}
	return response, err
}

// ListStorageAttachemnts returns list of attached storage blaocks to a worker node
func (r *worker) ListStorageAttachemnts(clusterIDOrName, workerID string, target ClusterTargetHeader) (VoulemeAttachments, error) {
	rawURL := fmt.Sprintf("/v2/storage/getAttachments?cluster=%s&worker=%s", clusterIDOrName, workerID)
	workerAttachements := VoulemeAttachments{}
	_, err := r.client.Get(rawURL, &workerAttachements, target.ToMap())
	if err != nil {
		return workerAttachements, err
	}
	return workerAttachements, err

}

func (r *worker) GetStorageAttachment(clusterIDOrName, workerID, volumeAttachmentID string, target ClusterTargetHeader) (VoulemeAttachment, error) {
	rawURL := fmt.Sprintf("/v2/storage/getAttachment?cluster=%s&worker=%s&volumeAttachmentID=%s", clusterIDOrName, workerID, volumeAttachmentID)
	workerAttachement := VoulemeAttachment{}
	_, err := r.client.Get(rawURL, &workerAttachement, target.ToMap())
	if err != nil {
		return workerAttachement, err
	}
	return workerAttachement, err

}

func (r *worker) CreateStorageAttachment(payload VolumeRequest, target ClusterTargetHeader) (VoulemeAttachment, error) {
	response := VoulemeAttachment{}
	_, err := r.client.Post("/v2/storage/createAttachment", payload, &response, target.ToMap())
	if err != nil {
		return response, err
	}
	return response, err
}

func (r *worker) DeleteStorageAttachment(payload VolumeRequest, target ClusterTargetHeader) (string, error) {
	var response string
	_, err := r.client.Post("/v2/storage/deleteAttachment", payload, &response, target.ToMap())
	if err != nil {
		return response, err
	}
	return response, err
}
