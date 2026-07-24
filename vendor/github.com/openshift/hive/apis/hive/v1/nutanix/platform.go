package nutanix

import (
	corev1 "k8s.io/api/core/v1"
)

// Platform stores any global configuration used for Nutanix platform
type Platform struct {
	// PrismCentral is the endpoint (address and port) to connect to the Prism Central.
	// This serves as the default Prism-Central.
	PrismCentral PrismEndpoint `json:"prismCentral"`

	// CredentialsSecretRef refers to a secret that contains the Nutanix account access
	// credentials.
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// CertificatesSecretRef refers to a secret that contains the Prism Central CA certificates
	// necessary for communicating with the Prism Central.
	// +optional
	CertificatesSecretRef corev1.LocalObjectReference `json:"certificatesSecretRef"`

	// FailureDomains configures failure domains for the Nutanix platform.
	// Required for using MachinePools
	FailureDomains []FailureDomain `json:"failureDomains,omitempty"`
}

// PrismEndpoint holds the endpoint address and port to access the Nutanix Prism Central or Element (cluster)
type PrismEndpoint struct {
	// address is the endpoint address (DNS name or IP address) of the Nutanix Prism Central or Element (cluster)
	Address string `json:"address"`

	// port is the port number to access the Nutanix Prism Central or Element (cluster)
	Port int32 `json:"port"`
}

// PrismElement holds the uuid, endpoint of the Prism Element (cluster)
type PrismElement struct {
	// UUID is the UUID of the Prism Element (cluster)
	UUID string `json:"uuid"`

	// Endpoint holds the address and port of the Prism Element
	// +optional
	Endpoint PrismEndpoint `json:"endpoint,omitempty"`

	// Name is the Prism Element (cluster) name.
	Name string `json:"name,omitempty"`
}

// FailureDomain configures failure domain information for the Nutanix platform.
type FailureDomain struct {
	// Name defines the unique name of a failure domain.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`^[0-9A-Za-z_.-@/]+$`
	Name string `json:"name"`

	// PrismElement holds the identification (name, UUID) and the optional endpoint address and
	// port of the Nutanix Prism Element. When a cluster-wide proxy is installed, this endpoint will,
	// by default, be accessed through the cluster-wide proxy configured for the platform.
	// If communication with this endpoint should bypass the proxy, add the endpoint to the install-config
	// `spec.noProxy` list in the proxy configuration.
	// +kubebuilder:validation:Required
	PrismElement PrismElement `json:"prismElement"`

	// SubnetUUIDs identifies the network subnets of the Prism Element.
	// Currently we only support one subnet for a failure domain.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +listType=set
	SubnetUUIDs []string `json:"subnetUUIDs"`

	// StorageContainers identifies the storage containers in the Prism Element.
	// +optional
	StorageContainers []StorageResourceReference `json:"storageContainers,omitempty"`

	// DataSourceImages identifies the datasource images in the Prism Element.
	// +optional
	DataSourceImages []StorageResourceReference `json:"dataSourceImages,omitempty"`
}

// StorageResourceReference holds reference information of a storage resource (storage container, data source image, etc.)
type StorageResourceReference struct {
	// ReferenceName is the identifier of the storage resource configured in the FailureDomain.
	// +optional
	ReferenceName string `json:"referenceName,omitempty"`

	// UUID is the UUID of the storage container resource in the Prism Element.
	// +kubebuilder:validation:Required
	UUID string `json:"uuid"`

	// Name is the name of the storage container resource in the Prism Element.
	// +optional
	Name string `json:"name,omitempty"`
}
