package openstack

import (
        "github.com/openshift/installer/pkg/ipnet"
)

type AciNetExtStruct struct {
        InfraVLAN               string          `json:"infraVlan,omitempty"`
        KubeApiVLAN             string          `json:"kubeApiVlan,omitempty"`
        ServiceVLAN             string          `json:"serviceVlan,omitempty"`
        Mtu                     string          `json:"mtu,omitempty"`
        ProvisionTar            string		`json:"provisionTar,omitempty"`
        NeutronCIDR             *ipnet.IPNet    `json:"neutronCIDR,omitempty"`
        InstallerHostSubnet	string          `json:"installerHostSubnet"`
}

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// Region specifies the OpenStack region where the cluster will be created.
	// Deprecated: this value is not used by the installer.
	Region string `json:"region"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on OpenStack for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Cloud is the name of OpenStack cloud to use from clouds.yaml.
	Cloud string `json:"cloud"`

	// ExternalNetwork is name of the external network in your OpenStack cluster.
	ExternalNetwork string `json:"externalNetwork"`

        // AciNetExt is the name of the map for APIC nested domain fields
        AciNetExt AciNetExtStruct `json:"aciNetExt"`

	// FlavorName is the name of the compute flavor to use for instances in this cluster.
	FlavorName string `json:"computeFlavor"`

	// LbFloatingIP is the IP address of an available floating IP in your OpenStack cluster
	// to associate with the OpenShift load balancer.
	LbFloatingIP string `json:"lbFloatingIP"`

	// ExternalDNS holds the IP addresses of dns servers that will
	// be added to the dns resolution of all instances in the cluster.
	// +optional
	ExternalDNS []string `json:"externalDNS"`

	// TrunkSupport holds a `0` or `1` value that indicates whether or not to use trunk ports
	// in your OpenShift cluster.
	TrunkSupport string `json:"trunkSupport"`

	// OctaviaSupport holds a `0` or `1` value that indicates whether your OpenStack
	// cluster supports Octavia Loadbalancing.
	OctaviaSupport string `json:"octaviaSupport"`

	// ClusterOSImage is either a URL to override the default OS image
	// for cluster nodes or an existing Glance image name.
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`
}
