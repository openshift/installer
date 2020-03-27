variable "tags" {
  type        = map(string)
  default     = {}
  description = "tags to be applied to created resources."
}

variable "cluster_id" {
  description = "The identifier for the cluster."
  type        = string
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

variable "external_lb_fqdn_v4" {
  description = "External API's LB fqdn for IPv4"
  type        = string
}

variable "external_lb_fqdn_v6" {
  description = "External API's LB fqdn for IPv6"
  type        = string
}

variable "internal_lb_ipaddress_v4" {
  description = "External API's LB IP v4 address"
  type        = string
}

variable "internal_lb_ipaddress_v6" {
  description = "External API's LB IP v6 address"
  type        = string
}

variable "virtual_network_id" {
  description = "The ID for Virtual Network that will be linked to the Private DNS zone."
  type        = string
}

variable "resource_group_name" {
  type        = string
  description = "Resource group for the deployment"
}

variable "private" {
  type        = bool
  description = "This value determines if this is a private cluster or not."
}

variable "use_ipv4" {
  type        = bool
  description = "This value determines if this is cluster should use IPv4 networking."
}

variable "use_ipv6" {
  type        = bool
  description = "This value determines if this is cluster should use IPv6 networking."
}

variable "emulate_single_stack_ipv6" {
  type        = bool
  description = "This determines whether a dual-stack cluster is configured to emulate single-stack IPv6."
}
