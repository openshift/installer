variable "tags" {
  type        = map(string)
  default     = {}
  description = "tags to be applied to created resources."
}

variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = string
}

variable "base_domain" {
  description = "The base domain used for public records"
  type        = string
}

variable "base_domain_resource_group_name" {
  description = "The resource group where the base domain is"
  type        = string
}

variable "external_lb_fqdn" {
  description = "External API's LB fqdn"
  type        = string
}

variable "internal_lb_ipaddress" {
  description = "External API's LB Ip address"
  type        = string
}

variable "private_dns_zone_name" {
  description = "private DNS zone name that should be used for records"
  type        = string
}

variable "etcd_count" {
  description = "The number of etcd members."
  type        = string
}

variable "etcd_ip_addresses" {
  description = "List of string IPs for machines running etcd members."
  type        = list(string)
  default     = []
}

variable "resource_group_name" {
  type        = string
  description = "Resource group for the deployment"
}

