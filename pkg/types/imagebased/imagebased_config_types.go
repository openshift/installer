package imagebased

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
)

const (
	// ImageBasedConfigVersion is the version supported by this package.
	ImageBasedConfigVersion = "v1beta1"
)

// Config is the API for specifying configuration for the image-based configuration ISO.
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Hostname is the desired hostname of the SNO node.
	Hostname string `json:"hostname,omitempty"`

	// NetworkConfig is a YAML manifest that can be processed by nmstate, using custom
	// marshaling/unmarshaling that will allow to populate nmstate config as plain yaml.
	NetworkConfig aiv1beta1.NetConfig `json:"networkConfig,omitempty"`

	// ReleaseRegistry is the container registry used to host the release image of the seed cluster.
	ReleaseRegistry string `json:"releaseRegistry,omitempty"`
}
