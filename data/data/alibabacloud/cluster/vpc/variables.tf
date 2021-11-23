variable "cluster_id" {
  type = string
}

variable "region_id" {
  type = string
}

variable "zone_ids" {
  type        = list(string)
  description = "The availability zones in which to create the masters."
}

variable "nat_gateway_zone_id" {
  type        = string
  description = "The availability zone in which to create the NAT gateway."
}

variable "vpc_cidr_block" {
  type = string
}

variable "resource_group_id" {
  type = string
}

variable "tags" {
  type        = map(string)
  description = "Tags to be applied to created resources."
}
