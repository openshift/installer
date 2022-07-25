variable "cluster_id" {
  type = string
}

variable "private_zone_id" {
  type = string
}

variable "resource_group_id" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = string
}

variable "base_domain" {
  description = "The base domain used for public records."
  type        = string
}

variable "slb_external_ip" {
  type        = string
  description = "External SLB IP address."
}

variable "slb_internal_ip" {
  type        = string
  description = "Internal SLB IP address."
}

variable "tags" {
  type        = map(string)
  description = "Tags to be applied to created resources."
}

variable "publish_strategy" {
  type        = string
  description = "The publishing strategy for endpoints like load balancers"
}
