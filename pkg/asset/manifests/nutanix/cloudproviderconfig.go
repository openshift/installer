package nutanix

import (
	"bytes"
	"encoding/json"
	"fmt"

	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
)

// CloudConfig is the config of Nutanix cloud provider
// ref: https://github.com/nutanix-cloud-native/cloud-provider-nutanix/blob/main/pkg/provider/config/config.go
type CloudConfig struct {
	PrismCentral         PrismEndpoint     `json:"prismCentral"`
	TopologyDiscovery    TopologyDiscovery `json:"topologyDiscovery"`
	EnableCustomLabeling bool              `json:"enableCustomLabeling"`
}

// TopologyDiscovery of the cloud provider.
type TopologyDiscovery struct {
	// Default type will be set to Prism via the newConfig function
	Type               TopologyDiscoveryType `json:"type"`
	TopologyCategories *TopologyCategories   `json:"topologyCategories"`
}

// TopologyDiscoveryType type alias.
type TopologyDiscoveryType string

const (
	// PrismTopologyDiscoveryType is the DiscoveryType for Prism provider.
	PrismTopologyDiscoveryType = TopologyDiscoveryType("Prism")
	// CategoriesTopologyDiscoveryType is the DiscoveryType for Categories provider.
	CategoriesTopologyDiscoveryType = TopologyDiscoveryType("Categories")
)

// TopologyInfo contains topology information.
type TopologyInfo struct {
	Zone   string `json:"zone"`
	Region string `json:"region"`
}

// TopologyCategories contains topology categories.
type TopologyCategories struct {
	ZoneCategory   string `json:"zoneCategory"`
	RegionCategory string `json:"regionCategory"`
}

// PrismEndpoint contains endpoint details for Prism provider.
type PrismEndpoint struct {
	// address is the endpoint address (DNS name or IP address) of the Nutanix Prism Central or Element (cluster)
	Address string `json:"address"`

	// port is the port number to access the Nutanix Prism Central or Element (cluster)
	Port int32 `json:"port"`

	// Pass credential information for the target Prism instance
	// +optional
	CredentialRef *CredentialReference `json:"credentialRef,omitempty"`
}

// CredentialKind type alias.
type CredentialKind string

// SecretKind a credential of type "Secret".
var SecretKind = CredentialKind("Secret")

// CredentialReference holds details of a credential.
type CredentialReference struct {
	// Kind of the Nutanix credential
	Kind CredentialKind `json:"kind"`

	// Name of the credential.
	Name string `json:"name"`
	// namespace of the credential.
	Namespace string `json:"namespace"`
}

// CloudConfigJSON returns the json string of the created CloudConfig
// based on the input nutanix platform config.
func CloudConfigJSON(nutanixPlatform *nutanixtypes.Platform) (string, error) {
	nutanixCloudConfig := CloudConfig{
		PrismCentral: PrismEndpoint{
			Address: nutanixPlatform.PrismCentral.Endpoint.Address,
			Port:    nutanixPlatform.PrismCentral.Endpoint.Port,
			CredentialRef: &CredentialReference{
				Kind:      "Secret",
				Name:      "nutanix-credentials",
				Namespace: "openshift-cloud-controller-manager",
			},
		},
		TopologyDiscovery: TopologyDiscovery{
			Type: PrismTopologyDiscoveryType,
		},
		EnableCustomLabeling: true,
	}
	configData, err := nutanixCloudConfig.JSON()
	if err != nil {
		return "", fmt.Errorf("could not create Nutanix Cloud provider config. %w", err)
	}

	return configData, nil
}

// JSON generates the cloud provider json config for the Nutanix platform.
func (config CloudConfig) JSON() (string, error) {
	buff := &bytes.Buffer{}
	encoder := json.NewEncoder(buff)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(config); err != nil {
		return "", err
	}

	return buff.String(), nil
}
