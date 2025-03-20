package containerv2

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

const (
	// EnableOutboundProtection configures a secure by default cluster to block all public outbound traffic
	EnableOutboundProtection = "enable-outbound-protection"
	// DisableOutboundProtection configures a secure by default cluster to allow all public outbound traffic
	DisableOutboundProtection = "disable-outbound-protection"
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

// VPCs interface
type VPCs interface {
	ListVPCs(target ClusterTargetHeader) ([]VPCConfig, error)
	SetOutboundTrafficProtection(string, bool, ClusterTargetHeader) error
	EnableSecureByDefault(string, bool, ClusterTargetHeader) error
}

func newVPCsAPI(c *client.Client) VPCs {
	return &vpc{
		client: c,
	}
}

// ListVPCs lists the vpcs
func (r *vpc) ListVPCs(target ClusterTargetHeader) ([]VPCConfig, error) {
	var successV []VPCConfig
	_, err := r.client.Get(fmt.Sprintf("/v2/vpc/getVPCs?provider=%s", target.Provider), &successV, target.ToMap())
	return successV, err
}

type OutboundTrafficProtectionRequest struct {
	Cluster   string `json:"cluster" binding:"required"`
	Operation string `json:"operation" binding:"required"`
}

// Set Outbound traffic protection
func (v *vpc) SetOutboundTrafficProtection(clusterID string, enable bool, target ClusterTargetHeader) error {
	request := OutboundTrafficProtectionRequest{
		Cluster:   clusterID,
		Operation: DisableOutboundProtection,
	}
	if enable {
		request.Operation = EnableOutboundProtection
	}

	_, err := v.client.Post("/network/v2/outbound-traffic-protection", request, nil, target.ToMap())

	return err
}

type EnableSecureByDefaultClusterRequest struct {
	Cluster                          string `json:"cluster" binding:"required"`
	DisableOutboundTrafficProtection bool   `json:"disableOutboundTrafficProtection,omitempty"`
}

// Enable Secure by Default
func (v *vpc) EnableSecureByDefault(clusterID string, enable bool, target ClusterTargetHeader) error {
	request := EnableSecureByDefaultClusterRequest{
		Cluster:                          clusterID,
		DisableOutboundTrafficProtection: enable,
	}

	_, err := v.client.Post("/network/v2/secure-by-default/enable", request, nil, target.ToMap())

	return err
}
