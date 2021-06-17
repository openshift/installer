variable "cluster_id" {
  type = string
}

variable "region_id" {
  type = string
}

variable "zone_ids" {
  type        = list(string)
  description = "The availability zones in which to create the masters and workers."
}

variable "resource_group_id" {
  type = string
}

variable "vpc_cidr_block" {
  type = string
}

variable "vswitch_cidr_blocks" {
  type        = list(string)
  description = "A list of IPv4 CIDRs."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Tags to be applied to created resources."
}
