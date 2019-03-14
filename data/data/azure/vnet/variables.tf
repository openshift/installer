variable "cidr_block" {
  type = "string"
}

variable "rg_name" {
  type = "string"
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
