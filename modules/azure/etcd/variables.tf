// Location is the Azure Location (East US, West US, etc)
variable "location" {
  type = "string"
}

variable "resource_group_name" {
  type = "string"
}

// Image refernce to use for master instances
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

// Count of etcd nodes to be created.
variable "etcd_count" {
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

variable "network_interface_ids" {
  type = "list"
}

variable "endpoints" {
  type = "list"
}
