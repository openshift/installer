variable "vnet_cidr" {
  type = "string"
}

variable "master_subnet_cidr" {
  type        = "string"
  description = "The subnet for the masters"
}

variable "node_subnet_cidr" {
  type        = "string"
  description = "The subnet for the workers"
}

variable "resource_group_name" {
  type        = "string"
  description = "Resource group for the deployment"
}

variable "cluster_id" {
  type = "string"
}

variable "private_master_endpoints" {
  description = "If set to true, private-facing ingress resources are created."
  default     = true
}

variable "public_master_endpoints" {
  description = "If set to true, public-facing ingress resources are created."
  default     = true
}

variable "region" {
  type        = "string"
  description = "The target Azure region for the cluster."
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "Azure tags to be applied to created resources."
}

variable "dns_label" {
  type        = "string"
  description = "The label used to build the dns name. i.e. <label>.<region>.cloudapp.azure.com"
}
