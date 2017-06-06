// Location is the Azure Location (East US, West US, etc)
variable "location" {
  type = "string"
}

variable "resource_group_name" {
  type = "string"
}

// Image refernce to use for worker instances
variable "image_reference" {
  type = "map"
}

// VM Size name
variable "vm_size" {
  type = "string"
}

// Storage account type
variable "storage_account_type" {
  type = "string"
}

// Count of worker nodes to be created.
variable "worker_count" {
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

variable "subnet" {
  type = "string"
}

variable "kube_image_url" {
  type = "string"
}

variable "kube_image_tag" {
  type = "string"
}

variable "kubeconfig_content" {
  type    = "string"
  default = ""
}

variable "tectonic_kube_dns_service_ip" {
  type = "string"
}

variable "cloud_provider" {
  type    = "string"
  default = "azure"
}

variable "kubelet_node_label" {
  type = "string"
}
