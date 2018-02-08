package terraformgenerator

import (
	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

// VMware defines all variables for this platform.
type VMware struct {
	ControllerDomain   string `json:"tectonic_vmware_controller_domain,omitempty"`
	EtcdClusters       string `json:"tectonic_vmware_etcd_clusters,omitempty"`
	EtcdDatacenters    string `json:"tectonic_vmware_etcd_datacenters,omitempty"`
	EtcdDatastores     string `json:"tectonic_vmware_etcd_datastores,omitempty"`
	EtcdGateways       string `json:"tectonic_vmware_etcd_gateways,omitempty"`
	EtcdHostnames      string `json:"tectonic_vmware_etcd_hostnames,omitempty"`
	EtcdIP             string `json:"tectonic_vmware_etcd_ip,omitempty"`
	EtcdMemory         string `json:"tectonic_vmware_etcd_memory,omitempty"`
	EtcdNetworks       string `json:"tectonic_vmware_etcd_networks,omitempty"`
	EtcdResourcePool   string `json:"tectonic_vmware_etcd_resource_pool,omitempty"`
	EtcdVCPU           string `json:"tectonic_vmware_etcd_vcpu,omitempty"`
	Folder             string `json:"tectonic_vmware_folder,omitempty"`
	IngressDomain      string `json:"tectonic_vmware_ingress_domain,omitempty"`
	MasterClusters     string `json:"tectonic_vmware_master_clusters,omitempty"`
	MasterDatacenters  string `json:"tectonic_vmware_master_datacenters,omitempty"`
	MasterDatastores   string `json:"tectonic_vmware_master_datastores,omitempty"`
	MasterGateways     string `json:"tectonic_vmware_master_gateways,omitempty"`
	MasterHostnames    string `json:"tectonic_vmware_master_hostnames,omitempty"`
	MasterIP           string `json:"tectonic_vmware_master_ip,omitempty"`
	MasterMemory       string `json:"tectonic_vmware_master_memory,omitempty"`
	MasterNetworks     string `json:"tectonic_vmware_master_networks,omitempty"`
	MasterResourcePool string `json:"tectonic_vmware_master_resource_pool,omitempty"`
	MasterVCPU         string `json:"tectonic_vmware_master_vcpu,omitempty"`
	NodeDNS            string `json:"tectonic_vmware_node_dns,omitempty"`
	Server             string `json:"tectonic_vmware_server,omitempty"`
	SSHAuthorizedKey   string `json:"tectonic_vmware_ssh_authorized_key,omitempty"`
	SSHPrivateKeyPath  string `json:"tectonic_vmware_ssh_private_key_path,omitempty"`
	SSLSelfSigned      string `json:"tectonic_vmware_sslselfsigned,omitempty"`
	Type               string `json:"tectonic_vmware_type,omitempty"`
	VMTemplate         string `json:"tectonic_vmware_vm_template,omitempty"`
	VMTemplateFolder   string `json:"tectonic_vmware_vm_template_folder,omitempty"`
	WorkerClusters     string `json:"tectonic_vmware_worker_clusters,omitempty"`
	WorkerDatacenters  string `json:"tectonic_vmware_worker_datacenters,omitempty"`
	WorkerDatastores   string `json:"tectonic_vmware_worker_datastores,omitempty"`
	WorkerGateways     string `json:"tectonic_vmware_worker_gateways,omitempty"`
	WorkerHostnames    string `json:"tectonic_vmware_worker_hostnames,omitempty"`
	WorkerIP           string `json:"tectonic_vmware_worker_ip,omitempty"`
	WorkerMemory       string `json:"tectonic_vmware_worker_memory,omitempty"`
	WorkerNetworks     string `json:"tectonic_vmware_worker_networks,omitempty"`
	WorkerResourcePool string `json:"tectonic_vmware_worker_resource_pool,omitempty"`
	WorkerVCPU         string `json:"tectonic_vmware_worker_vcpu,omitempty"`
}

// NewVMWare returns the config for VMware.
func NewVMWare(cluster config.Cluster) VMware {
	return VMware{
		ControllerDomain:   cluster.VMware.ControllerDomain,
		EtcdClusters:       cluster.VMware.Etcd.Clusters,
		EtcdDatacenters:    cluster.VMware.Etcd.Datacenters,
		EtcdDatastores:     cluster.VMware.Etcd.Datastores,
		EtcdGateways:       cluster.VMware.Etcd.Gateways,
		EtcdHostnames:      cluster.VMware.Etcd.Hostnames,
		EtcdIP:             cluster.VMware.Etcd.IP,
		EtcdMemory:         cluster.VMware.Etcd.Memory,
		EtcdNetworks:       cluster.VMware.Etcd.Networks,
		EtcdResourcePool:   cluster.VMware.Etcd.ResourcePool,
		EtcdVCPU:           cluster.VMware.Etcd.VCPU,
		Folder:             cluster.VMware.Folder,
		IngressDomain:      cluster.VMware.IngressDomain,
		MasterClusters:     cluster.VMware.Master.Clusters,
		MasterDatacenters:  cluster.VMware.Master.Datacenters,
		MasterDatastores:   cluster.VMware.Master.Datastores,
		MasterGateways:     cluster.VMware.Master.Gateways,
		MasterHostnames:    cluster.VMware.Master.Hostnames,
		MasterIP:           cluster.VMware.Master.IP,
		MasterMemory:       cluster.VMware.Master.Memory,
		MasterNetworks:     cluster.VMware.Master.Networks,
		MasterResourcePool: cluster.VMware.Master.ResourcePool,
		MasterVCPU:         cluster.VMware.Master.VCPU,
		NodeDNS:            cluster.VMware.NodeDNS,
		Server:             cluster.VMware.Server,
		SSHAuthorizedKey:   cluster.VMware.SSH.AuthorizedKey,
		SSHPrivateKeyPath:  cluster.VMware.SSH.PrivateKeyPath,
		SSLSelfSigned:      cluster.VMware.SSLSelfSigned,
		Type:               cluster.VMware.Type,
		VMTemplate:         cluster.VMware.VM.Template,
		VMTemplateFolder:   cluster.VMware.VM.TemplateFolder,
		WorkerClusters:     cluster.VMware.Worker.Clusters,
		WorkerDatacenters:  cluster.VMware.Worker.Datacenters,
		WorkerDatastores:   cluster.VMware.Worker.Datastores,
		WorkerGateways:     cluster.VMware.Worker.Gateways,
		WorkerHostnames:    cluster.VMware.Worker.Hostnames,
		WorkerIP:           cluster.VMware.Worker.IP,
		WorkerMemory:       cluster.VMware.Worker.Memory,
		WorkerNetworks:     cluster.VMware.Worker.Networks,
		WorkerResourcePool: cluster.VMware.Worker.ResourcePool,
		WorkerVCPU:         cluster.VMware.Worker.VCPU,
	}
}
