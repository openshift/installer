variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
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

variable "base_domain" {
  description = "The base domain used for public records."
  type        = string
}

variable "vpc_id" {
  description = "The VPC used to create the private route53 zone."
}

variable "cluster_id" {
  type        = string
  description = "The identifier for the cluster."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "AWS tags to be applied to created resources."
}

variable "api_external_lb_dns_name" {
  description = "External API's LB DNS name"
  type        = string
}

variable "api_external_lb_zone_id" {
  description = "External API's LB Zone ID"
  type        = string
}

variable "api_internal_lb_dns_name" {
  description = "Internal API's LB DNS name"
  type        = string
}

variable "api_internal_lb_zone_id" {
  description = "Internal API's LB Zone ID"
  type        = string
}

