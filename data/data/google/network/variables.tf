variable "bootstrap_instance" {
  type        = "string"
  description = "The bootstrap instance."
}

variable "bootstrap_instance_group" {
  type        = "string"
  description = "The instance group that hold the bootstrap instance in this region."
}

variable "cidr_block" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "extra_labels" {
  type        = "map"
  default     = {}
  description = "Extra GCP labels to be applied to created resources."
}

variable "master_instances" {
  type        = "list"
  description = "The master instances."
}

variable "master_instance_groups" {
  type        = "list"
  description = "The instance groups that hold the master instances in this region."
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
  description = "The target GCP region for the cluster."
}
