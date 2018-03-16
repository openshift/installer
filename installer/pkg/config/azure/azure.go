package azure

// Azure converts azure related config.
type Azure struct {
	CloudEnvironment string `json:"tectonic_azure_cloud_environment,omitempty" yaml:"cloudEnvironment,omitempty"`
	Etcd             `json:",inline" yaml:"etcd,omitempty"`
	External         `json:",inline" yaml:"external,omitempty"`
	ExtraTags        string `json:"tectonic_azure_extra_tags,omitempty" yaml:"extraTags,omitempty"`
	Master           `json:",inline" yaml:"master,omitempty"`
	PrivateCluster   string `json:"tectonic_azure_private_cluster,omitempty" yaml:"privateCluster,omitempty"`
	SSH              `json:",inline" yaml:"ssh,omitempty"`
	VNetCIDRBlock    string `json:"tectonic_azure_vnet_cidr_block,omitempty" yaml:"vNetCIDRBlock,omitempty"`
	Worker           `json:",inline" yaml:"worker,omitempty"`
}

// Etcd converts etcd related config.
type Etcd struct {
	StorageType string `json:"tectonic_azure_etcd_storage_type,omitempty" yaml:"storageType,omitempty"`
	VMSize      string `json:"tectonic_azure_etcd_vm_size,omitempty" yaml:"vmSize,omitempty"`
}

// External converts external related config.
type External struct {
	DNSZoneID      string `json:"tectonic_azure_external_dns_zone_id,omitempty" yaml:"dnsZoneID,omitempty"`
	MasterSubnetID string `json:"tectonic_azure_external_master_subnet_id,omitempty" yaml:"masterSubnetID,omitempty"`
	NSG            `json:",inline" yaml:"nsg,omitempty"`
	ResourceGroup  string `json:"tectonic_azure_external_resource_group,omitempty" yaml:"resourceGroup,omitempty"`
	VNetID         string `json:"tectonic_azure_external_vnet_id,omitempty" yaml:"vNetID,omitempty"`
	WorkerSubnetID string `json:"tectonic_azure_external_worker_subnet_id,omitempty" yaml:"workerSubnetID,omitempty"`
}

// Master converts master related config.
type Master struct {
	StorageType string `json:"tectonic_azure_master_storage_type,omitempty" yaml:"storageType,omitempty"`
	VMSize      string `json:"tectonic_azure_master_vm_size,omitempty" yaml:"vmSize,omitempty"`
}

// Network converts network related config.
type Network struct {
	External string `json:"tectonic_azure_ssh_network_external,omitempty" yaml:"external,omitempty"`
	Internal string `json:"tectonic_azure_ssh_network_internal,omitempty" yaml:"internal,omitempty"`
}

// NSG converts nsg related config.
type NSG struct {
	MasterID string `json:"tectonic_azure_external_nsg_master_id,omitempty" yaml:"masterID,omitempty"`
	WorkerID string `json:"tectonic_azure_external_nsg_worker_id,omitempty" yaml:"workerID,omitempty"`
}

// SSH converts ssh related config.
type SSH struct {
	Key     string `json:"tectonic_azure_ssh_key,omitempty" yaml:"key,omitempty"`
	Network `json:",inline" yaml:"network,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
	StorageType string `json:"tectonic_azure_worker_storage_type,omitempty" yaml:"storageType,omitempty"`
	VMSize      string `json:"tectonic_azure_worker_vm_size,omitempty" yaml:"vmSize,omitempty"`
}
