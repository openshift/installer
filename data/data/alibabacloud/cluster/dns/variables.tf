variable "cluster_id" {
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

variable "master_count" {
  type = number
}

variable "master_ips" {
  type = list(string)
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "Tags to be applied to created resources."
}
