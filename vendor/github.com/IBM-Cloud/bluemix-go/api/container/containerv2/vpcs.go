package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type VPCConfig struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Provider      string   `json:"provider"`
	ResourceGroup string   `json:"resourceGroup"`
	Zones         []string `json:"zones"`
}

type vpc struct {
	client *client.Client
}

//VPCs interface
type VPCs interface {
	ListVPCs(target ClusterTargetHeader) ([]VPCConfig, error)
}

func newVPCsAPI(c *client.Client) VPCs {
	return &vpc{
		client: c,
	}
}

//ListVPCs lists the vpcs
func (r *vpc) ListVPCs(target ClusterTargetHeader) ([]VPCConfig, error) {
	var successV []VPCConfig
	_, err := r.client.Get(fmt.Sprintf("/v2/vpc/getVPCs?provider=%s", target.Provider), &successV, target.ToMap())
	return successV, err
}
