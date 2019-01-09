variable "cidr_block" {
  type = "string"
}

variable "cluster_id" {
  type = "string"
}

variable "base_domain" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "extra_tags" {
  description = "Extra AWS tags to be applied to created resources."
  type        = "map"
  default     = {}
}

variable "private_controlplane_endpoints" {
  description = "If set to true, private-facing ingress resources are created."
  default     = true
}

variable "public_controlplane_endpoints" {
  description = "If set to true, public-facing ingress resources are created."
  default     = true
}

variable "region" {
  type        = "string"
  description = "The target AWS region for the cluster."
}
