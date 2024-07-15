package imagebased

import (
	"github.com/openshift/installer/pkg/types"
)

const (
	// SeedReconfigurationVersion is the current version of the
	// SeedReconfiguration struct.
	SeedReconfigurationVersion = 1

	// BlockDeviceLabel is the volume label to be used for the image-based
	// installer configuration ISO.
	BlockDeviceLabel = "cluster-config"
)

// SeedReconfiguration contains all the information that is required to
// transform a machine started from a single-node OpenShift (SNO) seed OCI image
// (which contains dummy seed configuration) into a SNO cluster with the desired
// configuration.
type SeedReconfiguration struct {
	// AdditionalTrustBundle keeps the PEM-encoded x.509 certificate bundle(s)
	// that will be added to the nodes' trusted certificate store.
	AdditionalTrustBundle AdditionalTrustBundle `json:"additionalTrustBundle,omitempty"`

	// APIVersion is the version of this struct and it is used to detect breaking
	// changes.
	APIVersion int `json:"api_version"`

	// BaseDomain is the desired base domain.
	BaseDomain string `json:"base_domain,omitempty"`

	// ClusterID is the desired cluster ID.
	ClusterID string `json:"cluster_id,omitempty"`

	// ClusterName is the desired cluster name.
	ClusterName string `json:"cluster_name,omitempty"`

	// ChronyConfig is the desired chrony configuration and it is used to populate
	// the /etc/chrony.conf on the node.
	ChronyConfig string `json:"chrony_config,omitempty"`

	// Hostname is the desired hostname of the node.
	Hostname string `json:"hostname,omitempty"`

	// InfraID is the desired infra ID.
	InfraID string `json:"infra_id,omitempty"`

	// KubeadminPasswordHash is the hash of the password for the kubeadmin
	// user, as can be found in the kubeadmin key of the kube-system/kubeadmin
	// secret. This will replace the kubeadmin password of the seed cluster.
	KubeadminPasswordHash string `json:"kubeadmin_password_hash,omitempty"`

	// KubeconfigCryptoRetention contains all the crypto material that is required
	// for the image-based installer to ensure that the generated kubeconfigs can
	// be used to access the cluster after its configuration.
	KubeconfigCryptoRetention KubeConfigCryptoRetention

	// MachineNetwork is the list of IP address pools for machines.
	// This field replaces MachineCIDR, and if set MachineCIDR must
	// be empty or match the first entry in the list.
	// Default is 10.0.0.0/16 for all platforms other than Power VS.
	// For Power VS, the default is 192.168.0.0/24.
	MachineNetwork string `json:"machine_network,omitempty"`

	// NodeIP is the desired IP address of the node.
	NodeIP string `json:"node_ip,omitempty"`

	// RawNMStateConfig contains the nmstate configuration YAML manifest as string.
	// Example nmstate configurations can be found here: https://nmstate.io/examples.html.
	RawNMStateConfig string `json:"raw_nm_state_config,omitempty"`

	// RelaseRegistry is the container registry that hosts the release image of
	// the seed cluster.
	ReleaseRegistry string `json:"release_registry,omitempty"`

	// SSHKey is the public Secure Shell (SSH) key that provides access to the
	// node.
	SSHKey string `json:"ssh_key,omitempty"`

	// Proxy defines the proxy settings for the cluster.
	// If unset, the cluster will not be configured to use a proxy.
	Proxy *types.Proxy `json:"proxy,omitempty"`

	// PullSecret is the secret to use when pulling images.
	PullSecret string `json:"pull_secret,omitempty"`
}

// KubeConfigCryptoRetention contains all the crypto material that is required
// for the image-based installer to ensure that the kubeconfigs can be used to
// access the cluster after its configuration.
type KubeConfigCryptoRetention struct {
	KubeAPICrypto KubeAPICrypto

	IngresssCrypto IngresssCrypto
}

// KubeAPICrypto contains the kubernetes API private keys and certificates that
// are used to generate and sign the cluster's cryptographic objects.
type KubeAPICrypto struct {
	ServingCrypto ServingCrypto

	ClientAuthCrypto ClientAuthCrypto
}

// ServingCrypto contains the kubernetes API private keys that are used to
// generate the cluster's certificates.
type ServingCrypto struct {
	// LocalhostSignerPrivateKey is a PEM-encoded X.509 key.
	LocalhostSignerPrivateKey string `json:"localhost_signer_private_key,omitempty"`

	// ServiceNetworkSignerPrivateKey is a PEM-encoded X.509 key.
	ServiceNetworkSignerPrivateKey string `json:"service_network_signer_private_key,omitempty"`

	// LoadbalancerSignerPrivateKey is a PEM-encoded X.509 key.
	LoadbalancerSignerPrivateKey string `json:"loadbalancer_external_signer_private_key,omitempty"`
}

// ClientAuthCrypto contains the CA certificate used to sign the cluster's
// cryptographic objects.
type ClientAuthCrypto struct {
	// AdminCACertificate is a PEM-encoded X.509 certificate.
	AdminCACertificate string `json:"admin_ca_certificate,omitempty"`
}

// IngresssCrypto contains the ingrees CA certificate.
type IngresssCrypto struct {
	// IngressCA is a PEM-encoded X.509 certificate.
	IngressCA string `json:"ingress_ca,omitempty"`
}

// AdditionalTrustBundle represents the PEM-encoded X.509 certificate bundle
// that will be added to the nodes' trusted certificate store.
type AdditionalTrustBundle struct {
	// UserCaBundle keeps the contents of the user-ca-bundle ConfigMap in the
	// openshift-config namepace.
	UserCaBundle string `json:"userCaBundle"`

	// ProxyConfigmapName is the Proxy CR trustedCA ConfigMap name.
	ProxyConfigmapName string `json:"proxyConfigmapName"`

	// ProxyConfigampBundle keeps the contents of the ProxyConfigmapName ConfigMap.
	// It must be equal to the UserCaBundle when  ProxyConfigmapName is
	// user-ca-bundle.
	ProxyConfigmapBundle string `json:"proxyConfigmapBundle"`
}
