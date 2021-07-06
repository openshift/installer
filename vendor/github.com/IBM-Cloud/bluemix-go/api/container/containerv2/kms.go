package containerv2

import (
	"github.com/IBM-Cloud/bluemix-go/client"
)

const (
	account       = "X-Auth-Resource-Account"
	resourceGroup = "X-Auth-Resource-Group"
)

//Request body to attach a KMS to a cluster
type KmsEnableReq struct {
	Cluster         string `json:"cluster"`
	Kms             string `json:"instance_id"`
	Crk             string `json:"crk_id"`
	PrivateEndpoint bool   `json:"private_endpoint"`
}

//ClusterHeader ...
type ClusterHeader struct {
	AccountID     string
	ResourceGroup string
}

//CreateMap ...
func (c ClusterHeader) CreateMap() map[string]string {
	m := make(map[string]string, 3)
	m[account] = c.AccountID
	m[resourceGroup] = c.ResourceGroup
	return m
}

type kms struct {
	client *client.Client
}

//Kms interface
type Kms interface {
	EnableKms(enableKmsReq KmsEnableReq, target ClusterHeader) error
}

func newKmsAPI(c *client.Client) Kms {
	return &kms{
		client: c,
	}
}

func (r *kms) EnableKms(enableKmsReq KmsEnableReq, target ClusterHeader) error {
	_, err := r.client.Post("/v2/enableKMS", enableKmsReq, nil, target.CreateMap())
	return err
}
