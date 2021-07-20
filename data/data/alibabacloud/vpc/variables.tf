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

variable "vpc_cidr_block" {
  type = string
}
 
variable "resource_group_id" {
  type = string
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Tags to be applied to created resources."
}
