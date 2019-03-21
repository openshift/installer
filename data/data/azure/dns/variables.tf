variable "tags" {
  type        = "map"
  default     = {}
  description = "AWS tags to be applied to created resources."
}

variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = "string"
}

variable "cluster_name" {
  description = "The name for the cluster without the random suffix"
  type        = "string"
}

variable "base_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = "string"
}

variable "external_lb_dns_label" {
  description = "External API's LB DNS name"
  type        = "string"
}

variable "internal_lb_ipaddress" {
  description = "External API's LB Ip address"
  type        = "string"
}

variable "internal_dns_resolution_vnet_id" {
  description = "the vnet id to be attached to the private DNS zone"
  type        = "string"
}

variable "resource_group_name" {
  type = "string"
  description = "Resource group for the deployment"
}

