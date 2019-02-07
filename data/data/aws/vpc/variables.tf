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

variable "region" {
  type        = "string"
  description = "The target AWS region for the cluster."
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "AWS tags to be applied to created resources."
}
