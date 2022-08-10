package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

// GetDedicatedHostFlavor is a response to a dedicated host pool get request.
// swagger:model
type GetDedicatedHostFlavor struct {
	ID              string            `json:"id"`
	FlavorClass     string            `json:"flavorClass"`
	Region          string            `json:"region"`
	Zone            string            `json:"zone"`
	Deprecated      bool              `json:"deprecated"`
	MaxVCPUs        int               `json:"maxVCPUs"`
	MaxMemory       int               `json:"maxMemory"`
	InstanceStorage []InstanceStorage `json:"instanceStorage"`
}

// GetDedicatedHostFlavors is a response to a dedicated host pool list request.
// swagger:model
type GetDedicatedHostFlavors []GetDedicatedHostFlavor

// InstanceStorage type for describing an instance disk configuration
// swagger:model
type InstanceStorage struct {
	Count int `json:"count"`
	// the size of each individual device in GB
	Size int `json:"size"`
}

//DedicatedHostFlavor ...
type DedicatedHostFlavor interface {
	ListDedicatedHostFlavors(zone string, target ClusterTargetHeader) (GetDedicatedHostFlavors, error)
}

type dedicatedhostflavor struct {
	client *client.Client
}

func newDedicatedHostFlavorAPI(c *client.Client) DedicatedHostFlavor {
	return &dedicatedhostflavor{
		client: c,
	}
}

// GetDedicatedHostFlavor calls the API to list dedicated host s
func (w *dedicatedhostflavor) ListDedicatedHostFlavors(zone string, target ClusterTargetHeader) (GetDedicatedHostFlavors, error) {
	successV := GetDedicatedHostFlavors{}
	_, err := w.client.Get(fmt.Sprintf("/v2/getDedicatedHostFlavors?provider=vpc-gen2&zone=%s", zone), &successV, target.ToMap())
	return successV, err
}
