package hpcs

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type EndpointResp struct {
	InstanceID string    `json:"instance_id"`
	Kms        Endpoints `json:"kms"`
	Ep11       Endpoints `json:"ep11"`
}
type Endpoints struct {
	Public  string `json:"public"`
	Private string `json:"private"`
}

type EndpointRepository interface {
	GetAPIEndpoint(instanceID string) (EndpointResp, error)
}

type hpcsRepository struct {
	client *client.Client
}

func NewHpcsEndpointRepository(c *client.Client) EndpointRepository {
	return &hpcsRepository{
		client: c,
	}
}

func (r *hpcsRepository) GetAPIEndpoint(instanceID string) (EndpointResp, error) {
	res := EndpointResp{}
	_, err := r.client.Get(fmt.Sprintf("/instances/%s", instanceID), &res)
	if err != nil {
		return EndpointResp{}, err
	}
	return res, nil
}
