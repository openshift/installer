variable "vnet_cidr" {
  type = string
}

variable "resource_group_name" {
  type        = string
  description = "Resource group for the deployment"
}

variable "cluster_id" {
  type = string
}

variable "region" {
  type        = string
  description = "The target Azure region for the cluster."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Azure tags to be applied to created resources."
}

variable "dns_label" {
  type        = string
  description = "The label used to build the dns name. i.e. <label>.<region>.cloudapp.azure.com"
}

variable "preexisting_network" {
  type        = bool
  description = "This value determines if a vnet already exists or not. If true, then will not create a new vnet, subnet, or nsg's"
  default     = false
}

variable "network_resource_group_name" {
  type        = string
  description = "This is the name of the network resource group for new or existing network resources"
}

variable "virtual_network_name" {
  type        = string
  description = "This is the name of the virtual network, new or existing"
}

variable "master_subnet" {
  type        = string
  description = "This is the name of the subnet used for the control plane, new or existing"
}

variable "worker_subnet" {
  type        = string
  description = "This is the name of the subnet used for the compute nodes, new or existing"
}

variable "private" {
  type        = bool
  description = "The determines if this is a private/internal cluster or not."
}
