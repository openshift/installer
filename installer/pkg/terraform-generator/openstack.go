package terraformgenerator

import (
	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

// OpenStack defines all variables for this platform.
type OpenStack struct {
	DisableFloatingIP string `json:"tectonic_openstack_disable_floatingip,omitempty"`
	DNSNameservers    string `json:"tectonic_openstack_dns_nameservers,omitempty"`
	EtcdFlavorID      string `json:"tectonic_openstack_etcd_flavor_id,omitempty"`
	EtcdFlavorName    string `json:"tectonic_openstack_etcd_flavor_name,omitempty"`
	ExternalGatewayID string `json:"tectonic_openstack_external_gateway_id,omitempty"`
	FloatingIPPool    string `json:"tectonic_openstack_floatingip_pool,omitempty"`
	ImageID           string `json:"tectonic_openstack_image_id,omitempty"`
	ImageName         string `json:"tectonic_openstack_image_name,omitempty"`
	LBProvider        string `json:"tectonic_openstack_lb_provider,omitempty"`
	MasterFlavorID    string `json:"tectonic_openstack_master_flavor_id,omitempty"`
	MasterFlavorName  string `json:"tectonic_openstack_master_flavor_name,omitempty"`
	SubnetCIDR        string `json:"tectonic_openstack_subnet_cidr,omitempty"`
	WorkerFlavorID    string `json:"tectonic_openstack_worker_flavor_id,omitempty"`
	WorkerFlavorName  string `json:"tectonic_openstack_worker_flavor_name,omitempty"`
}

// NewOpenStack returns the config for OpenStack.
func NewOpenStack(cluster config.Cluster) OpenStack {
	return OpenStack{
		DisableFloatingIP: cluster.OpenStack.DisableFloatingIP,
		DNSNameservers:    cluster.OpenStack.DNSNameservers,
		EtcdFlavorID:      cluster.OpenStack.EtcdFlavor.ID,
		EtcdFlavorName:    cluster.OpenStack.EtcdFlavor.Name,
		ExternalGatewayID: cluster.OpenStack.ExternalGatewayID,
		FloatingIPPool:    cluster.OpenStack.FloatingIPPool,
		ImageID:           cluster.OpenStack.Image.ID,
		ImageName:         cluster.OpenStack.Image.Name,
		LBProvider:        cluster.OpenStack.LBProvider,
		MasterFlavorID:    cluster.OpenStack.MasterFlavor.ID,
		MasterFlavorName:  cluster.OpenStack.MasterFlavor.Name,
		SubnetCIDR:        cluster.OpenStack.SubnetCIDR,
		WorkerFlavorID:    cluster.OpenStack.WorkerFlavor.ID,
		WorkerFlavorName:  cluster.OpenStack.WorkerFlavor.Name,
	}
}
