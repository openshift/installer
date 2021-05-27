package containerv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

//AddOn ...
type AddOn struct {
	AllowedUpgradeVersion []string    `json:"allowed_upgrade_versions,omitempty"`
	Deprecated            bool        `json:"deprecated"`
	HealthState           string      `json:"healthState,omitempty"`
	HealthStatus          string      `json:"healthStatus,omitempty"`
	MinKubeVersion        string      `json:"minKubeVersion,omitempty"`
	MinOCPVersion         string      `json:"minOCPVersion,omitempty"`
	Name                  string      `json:"name"`
	Options               interface{} `json:"options,omitempty"`
	SupportedKubeRange    string      `json:"supportedKubeRange,omitempty"`
	TargetVersion         string      `json:"targetVersion,omitempty"`
	Version               string      `json:"version,omitempty"`
	VlanSpanningRequired  bool        `json:"vlan_spanning_required"`
}

//GetAddOns ...
type GetAddOns struct {
	AddonsList []AddOn `json:"addons"`
}

//ConfigureAddOns ...
type ConfigureAddOns struct {
	AddonsList []AddOn `json:"addons"`
	Enable     bool    `json:"enable"`
	Update     bool    `json:"update"`
}

// AddOnsResponse ...
type AddOnsResponse struct {
	MissingDeps    interface{} `json:"missingDeps,omitempty"`
	OrphanedAddons interface{} `json:"orphanedAddons,omitempty"`
}

//AddOns ...
type AddOns interface {
	GetAddons(clusterName string, target ClusterTargetHeader) ([]AddOn, error)
	ConfigureAddons(clusterName string, params *ConfigureAddOns, target ClusterTargetHeader) (AddOnsResponse, error)
}

type addons struct {
	client *client.Client
}

func newAddOnsAPI(c *client.Client) AddOns {
	return &addons{
		client: c,
	}
}

//GetAddon ...
func (r *addons) GetAddons(clusterName string, target ClusterTargetHeader) ([]AddOn, error) {
	rawURL := fmt.Sprintf("/v1/clusters/%s/addons", clusterName)
	addonsList := GetAddOns{}
	_, err := r.client.Get(rawURL, &addonsList.AddonsList, target.ToMap())
	if err != nil {
		return addonsList.AddonsList, err
	}

	return addonsList.AddonsList, err
}

// ConfigureAddon ...
func (r *addons) ConfigureAddons(clusterName string, params *ConfigureAddOns, target ClusterTargetHeader) (AddOnsResponse, error) {
	rawURL := fmt.Sprintf("/v1/clusters/%s/addons", clusterName)
	resp := AddOnsResponse{}
	_, err := r.client.Patch(rawURL, params, &resp, target.ToMap())
	return resp, err
}
