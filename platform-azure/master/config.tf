// Location is the Azure Location (East US, West US, etc)
variable "location" {
  type    = "string"
  default = "East US"
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

// Kubelet version to run on the master.
variable "kubelet_version" {
  type = "string"
}

// Count of master nodes to be created.
variable "count" {
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

variable "ssh_key" {
  type = "string"
}
