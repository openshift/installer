variable "tectonic_azure_region" {
  type = "string"
}

// The amount of etcd nodes to be created.
// Example: `1`
variable "tectonic_etcd_count" {
  type    = "string"
  default = "1"
}

// The name of the cluster.
// This will be prepended to `tectonic_base_domain` resulting in the URL to the Tectonic console.
// Example: `demo`
variable "tectonic_cluster_name" {
  type = "string"
}

// The flavor ID as given in `azure flavor list`.
// Specifies the size (CPU/Memory/Drive) of the VM.
variable "tectonic_azure_vm_size" {
  type    = "string"
  default = "Standard_D2_v2"
}

variable "tectonic_azure_resource_group_name" {}

variable "tectonic_azure_dns_zone_name" {
  default = "azure.ifup.org"
}
