variable "az_count" {
  type = "string"
}

variable "cidr_block" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "external_vpc_id" {
  type = "string"
}

variable "external_master_subnets" {
  type = "list"
}

variable "external_worker_subnets" {
  type = "list"
}

variable "extra_tags" {
  description = "Extra AWS tags to be applied to created resources."
  type        = "map"
  default     = {}
}
