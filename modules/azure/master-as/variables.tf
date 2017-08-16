// Location is the Azure Location (East US, West US, etc)
variable "location" {
  type = "string"
}

variable "resource_group_name" {
  type = "string"
}

variable "cluster_id" {
  type = "string"
}

// VM Size name
variable "vm_size" {
  type = "string"
}

// Storage account type
variable "storage_type" {
  type = "string"
}

variable "storage_id" {
  type = "string"
}

// The base DNS domain of the cluster.
// Example: `azure.dev.coreos.systems`
variable "base_domain" {
  type = "string"
}

// The name of the cluster.
variable "cluster_name" {
  type = "string"
}

variable "public_ssh_key" {
  type = "string"
}

variable "virtual_network" {
  type = "string"
}

variable "kube_image_url" {
  type = "string"
}

variable "kube_image_tag" {
  type = "string"
}

variable "kubeconfig_content" {
  type = "string"
}

// Count of master nodes to be created.
variable "master_count" {
  type = "string"
}

variable "tectonic_kube_dns_service_ip" {
  type = "string"
}

variable "cloud_provider" {
  type    = "string"
  default = "azure"
}

variable "cloud_provider_config" {
  description = "Content of cloud provider config"
  type        = "string"
}

variable "kubelet_node_label" {
  type = "string"
}

variable "kubelet_node_taints" {
  type = "string"
}

variable "kubelet_cni_bin_dir" {
  type = "string"
}

variable "bootkube_service" {
  type        = "string"
  description = "The content of the bootkube systemd service unit"
}

variable "tectonic_service" {
  type        = "string"
  description = "The content of the tectonic installer systemd service unit"
}

variable "tectonic_service_disabled" {
  description = "Specifies whether the tectonic installer systemd unit will be disabled. If true, no tectonic assets will be deployed"
  default     = false
}

variable "network_interface_ids" {
  type        = "list"
  description = "List of NICs to use for master VMs"
}

variable "versions" {
  description = "(internal) Versions of the components to use"
  type        = "map"
}

variable "cl_channel" {
  type = "string"
}

variable "extra_tags" {
  type = "map"
}
