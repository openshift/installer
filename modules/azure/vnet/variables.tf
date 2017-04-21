variable "tectonic_azure_external_vnet_id" {
  type    = "string"
  default = ""
}

variable "tectonic_azure_vnet_cidr_block" {
  type    = "string"
  default = "10.0.0.0/16"
}

variable "tectonic_cluster_name" {
  type = "string"
}

variable "tectonic_azure_external_vnet_master_subnets" {
  type    = "list"
  default = []
}

variable "tectonic_azure_external_vnet_worker_subnets" {
  type    = "list"
  default = []
}

variable "resource_group_name" {
  type = "string"
}

variable "vnet_cidr_block" {
  type = "string"
}

variable "location" {
  type = "string"
}

variable "external_vnet_name" {
  type    = "string"
  default = ""
}
