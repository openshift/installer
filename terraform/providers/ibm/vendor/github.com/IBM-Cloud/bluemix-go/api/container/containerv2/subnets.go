package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type SubnetConfig struct {
	AvailableIPv4AddressCount int    `json:"availableIPv4AddressCount"`
	ID                        string `json:"id"`
	Ipv4CIDRBlock             string `json:"ipv4CIDRBlock"`
	Name                      string `json:"name"`
	PublicGatewayID           string `json:"publicGatewayID"`
	PublicGatewayName         string `json:"publicGatewayName"`
	VpcID                     string `json:"vpcID"`
	VpcName                   string `json:"vpcName"`
	Zone                      string `json:"zone"`
}

type subnet struct {
	client *client.Client
}

//Subnets interface
type Subnets interface {
	ListSubnets(vpcID, zone string, target ClusterTargetHeader) ([]SubnetConfig, error)
}

func newSubnetsAPI(c *client.Client) Subnets {
	return &subnet{
		client: c,
	}
}

//ListSubnets list the subnets for a given VPC
func (r *subnet) ListSubnets(vpcID, zone string, target ClusterTargetHeader) ([]SubnetConfig, error) {
	var successV []SubnetConfig
	_, err := r.client.Get(fmt.Sprintf("/v2/vpc/getSubnets?vpc=%s&provider=%s&zone=%s", vpcID, target.Provider, zone), &successV, target.ToMap())
	return successV, err
}
