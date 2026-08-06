package gcp

import (
	corev1 "k8s.io/api/core/v1"
)

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// CredentialsSecretRef refers to a secret that contains the GCP account access credentials.
	// +optional
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// Region specifies the GCP region where the cluster will be created.
	Region string `json:"region"`

	// PrivateSericeConnect allows users to enable access to the cluster's API server using GCP
	// Private Service Connect. It includes a forwarding rule paired with a Service Attachment
	// across GCP accounts and allows clients to connect to services using GCP internal networking
	// of using public load balancers.
	// +optional
	PrivateServiceConnect *PrivateServiceConnect `json:"privateServiceConnect,omitempty"`

	// DiscardLocalSsdOnHibernate passes the specified value through to the GCP API to indicate
	// whether the content of any local SSDs should be preserved or discarded. See
	// https://cloud.google.com/compute/docs/disks/local-ssd#stop_instance
	// This field is required when attempting to hibernate clusters with instances possessing
	// SSDs -- e.g. those with GPUs.
	// +optional
	DiscardLocalSsdOnHibernate *bool `json:"discardLocalSsdOnHibernate,omitempty"`
}

// PrivateServiceConnectAccess configures access to the cluster API using GCP Private Service Connect
type PrivateServiceConnect struct {
	// Enabled specifies if Private Service Connect is to be enabled on the cluster.
	Enabled bool `json:"enabled"`

	// ServiceAttachment configures the service attachment to be used by the cluster.
	// +optional
	ServiceAttachment *ServiceAttachment `json:"serviceAttachment,omitempty"`
}

// ServiceAttachment configures the service attachment to be used by the cluster
type ServiceAttachment struct {
	// Subnet configures the subnetwork that contains the service attachment.
	// +optional
	Subnet *ServiceAttachmentSubnet `json:"subnet,omitempty"`
}

// ServiceAttachmentSubnet configures the subnetwork used by the service attachment
type ServiceAttachmentSubnet struct {
	// Cidr specifies the cidr to use when creating a service attachment subnet.
	// +optional
	Cidr string `json:"cidr,omitempty"`

	// Existing specifies a pre-existing subnet to use instead of creating a new service attachment subnet.
	// This is required when using BYO VPCs. It must be in the same region as the api-int load balancer, be
	// configured with a purpose of "Private Service Connect", and have sufficient routing and firewall rules
	// to access the api-int load balancer.
	// +optional
	Existing *ServiceAttachmentSubnetExisting `json:"existing,omitempty"`
}

// ServiceAttachmentSubnetExisting describes the existing subnet.
type ServiceAttachmentSubnetExisting struct {
	// Name specifies the name of the existing subnet.
	Name string `json:"name"`

	// Project specifies the project the subnet exists in.
	// This is required for Shared VPC.
	// +optional
	Project string `json:"project,omitempty"`
}

// PlatformStatus contains the observed state on GCP platform.
type PlatformStatus struct {
	// PrivateServiceConnect contains the private service connect resource references
	// +optional
	PrivateServiceConnect *PrivateServiceConnectStatus `json:"privateServiceConnect,omitempty"`
}

// PrivateServiceConnectStatus contains the observed state for PrivateServiceConnect resources.
type PrivateServiceConnectStatus struct {
	// Endpoint is the selfLink of the endpoint created for the cluster.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// EndpointAddress is the selfLink of the address created for the cluster endpoint.
	// +optional
	EndpointAddress string `json:"endpointAddress,omitempty"`

	// ServiceAttachment is the selfLink of the service attachment created for the clsuter.
	// +optional
	ServiceAttachment string `json:"serviceAttachment,omitempty"`

	// ServiceAttachmentFirewall is the selfLink of the firewall that allows traffic between
	// the service attachment and the cluster's internal api load balancer.
	// +optional
	ServiceAttachmentFirewall string `json:"serviceAttachmentFirewall,omitempty"`

	// ServiceAttachmentSubnet is the selfLink of the subnet that will contain the service attachment.
	// +optional
	ServiceAttachmentSubnet string `json:"serviceAttachmentSubnet,omitempty"`
}
