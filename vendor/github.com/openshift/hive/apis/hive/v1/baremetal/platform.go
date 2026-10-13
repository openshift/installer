package baremetal

import corev1 "k8s.io/api/core/v1"

// Platform stores the global configuration for the cluster.
type Platform struct {
	// LibvirtSSHPrivateKeySecretRef is the reference to the secret that contains the private SSH key to use
	// for access to the libvirt provisioning host.
	// The SSH private key is expected to be in the secret data under the "ssh-privatekey" key.
	LibvirtSSHPrivateKeySecretRef corev1.LocalObjectReference `json:"libvirtSSHPrivateKeySecretRef"`
}
