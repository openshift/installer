variable "cidr_block" {
  type = "string"
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
  description = "The target AWS region for the cluster."
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "AWS tags to be applied to created resources."
}
