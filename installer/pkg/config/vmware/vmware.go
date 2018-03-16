package vmware

// Etcd converts etcd related config.
type Etcd struct {
	Clusters     string `json:"tectonic_vmware_etcd_clusters,omitempty" yaml:"clusters,omitempty"`
	Datacenters  string `json:"tectonic_vmware_etcd_datacenters,omitempty" yaml:"datacenters,omitempty"`
	Datastores   string `json:"tectonic_vmware_etcd_datastores,omitempty" yaml:"datastores,omitempty"`
	Gateways     string `json:"tectonic_vmware_etcd_gateways,omitempty" yaml:"gateways,omitempty"`
	Hostnames    string `json:"tectonic_vmware_etcd_hostnames,omitempty" yaml:"hostnames,omitempty"`
	IP           string `json:"tectonic_vmware_etcd_ip,omitempty" yaml:"ip,omitempty"`
	Memory       string `json:"tectonic_vmware_etcd_memory,omitempty" yaml:"memory,omitempty"`
	Networks     string `json:"tectonic_vmware_etcd_networks,omitempty" yaml:"networks,omitempty"`
	ResourcePool string `json:"tectonic_vmware_etcd_resource_pool,omitempty" yaml:"resourcePool,omitempty"`
	VCPU         string `json:"tectonic_vmware_etcd_vcpu,omitempty" yaml:"vCPU,omitempty"`
}

// Master converts master related config.
type Master struct {
	Clusters     string `json:"tectonic_vmware_master_clusters,omitempty" yaml:"clusters,omitempty"`
	Datacenters  string `json:"tectonic_vmware_master_datacenters,omitempty" yaml:"datacenters,omitempty"`
	Datastores   string `json:"tectonic_vmware_master_datastores,omitempty" yaml:"datastores,omitempty"`
	Gateways     string `json:"tectonic_vmware_master_gateways,omitempty" yaml:"gateways,omitempty"`
	Hostnames    string `json:"tectonic_vmware_master_hostnames,omitempty" yaml:"hostnames,omitempty"`
	IP           string `json:"tectonic_vmware_master_ip,omitempty" yaml:"ip,omitempty"`
	Memory       string `json:"tectonic_vmware_master_memory,omitempty" yaml:"memory,omitempty"`
	Networks     string `json:"tectonic_vmware_master_networks,omitempty" yaml:"networks,omitempty"`
	ResourcePool string `json:"tectonic_vmware_master_resource_pool,omitempty" yaml:"resourcePool,omitempty"`
	VCPU         string `json:"tectonic_vmware_master_vcpu,omitempty" yaml:"vCPU,omitempty"`
}

// SSH converts ssh related config.
type SSH struct {
	AuthorizedKey  string `json:"tectonic_vmware_ssh_authorized_key,omitempty" yaml:"authorizedKey,omitempty"`
	PrivateKeyPath string `json:"tectonic_vmware_ssh_private_key_path,omitempty" yaml:"privateKeyPath,omitempty"`
}

// VM converts vm related config.
type VM struct {
	Template       string `json:"tectonic_vmware_vm_template,omitempty" yaml:"template,omitempty"`
	TemplateFolder string `json:"tectonic_vmware_vm_template_folder,omitempty" yaml:"templateFolder,omitempty"`
}

// VMware converts VMware related config.
type VMware struct {
	ControllerDomain string `json:"tectonic_vmware_controller_domain,omitempty" yaml:"controllerDomain,omitempty"`
	Etcd             `json:",inline" yaml:"etcd,omitempty"`
	Folder           string `json:"tectonic_vmware_folder,omitempty" yaml:"folder,omitempty"`
	IngressDomain    string `json:"tectonic_vmware_ingress_domain,omitempty" yaml:"ingressDomain,omitempty"`
	Master           `json:",inline" yaml:"master,omitempty"`
	NodeDNS          string `json:"tectonic_vmware_node_dns,omitempty" yaml:"nodeDNS,omitempty"`
	Server           string `json:"tectonic_vmware_server,omitempty" yaml:"server,omitempty"`
	SSH              `json:",inline" yaml:"ssh,omitempty"`
	SSLSelfSigned    string `json:"tectonic_vmware_sslselfsigned,omitempty" yaml:"sslSelfsigned,omitempty"`
	Type             string `json:"tectonic_vmware_type,omitempty" yaml:"type,omitempty"`
	VM               `json:",inline" yaml:"vm,omitempty"`
	Worker           `json:",inline" yaml:"worker,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
	Clusters     string `json:"tectonic_vmware_worker_clusters,omitempty" yaml:"clusters,omitempty"`
	Datacenters  string `json:"tectonic_vmware_worker_datacenters,omitempty" yaml:"datacenters,omitempty"`
	Datastores   string `json:"tectonic_vmware_worker_datastores,omitempty" yaml:"datastores,omitempty"`
	Gateways     string `json:"tectonic_vmware_worker_gateways,omitempty" yaml:"gateways,omitempty"`
	Hostnames    string `json:"tectonic_vmware_worker_hostnames,omitempty" yaml:"hostnames,omitempty"`
	IP           string `json:"tectonic_vmware_worker_ip,omitempty" yaml:"ip,omitempty"`
	Memory       string `json:"tectonic_vmware_worker_memory,omitempty" yaml:"memory,omitempty"`
	Networks     string `json:"tectonic_vmware_worker_networks,omitempty" yaml:"networks,omitempty"`
	ResourcePool string `json:"tectonic_vmware_worker_resource_pool,omitempty" yaml:"resourcePool,omitempty"`
	VCPU         string `json:"tectonic_vmware_worker_vcpu,omitempty" yaml:"vCPU,omitempty"`
}
