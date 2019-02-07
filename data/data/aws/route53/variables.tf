variable "cluster_name" {
  description = "The name of the cluster"
  type        = "string"
}

variable "base_domain" {
  description = "The base domain used in records"
  type        = "string"
}

// AWS specific internal zone variables

variable "private_zone_id" {
  description = "Route53 Private Zone ID"
  type        = "string"
}

variable "api_external_lb_dns_name" {
  description = "External API's LB DNS name"
  type        = "string"
}

variable "api_external_lb_zone_id" {
  description = "External API's LB Zone ID"
  type        = "string"
}

variable "api_internal_lb_dns_name" {
  description = "Internal API's LB DNS name"
  type        = "string"
}

variable "api_internal_lb_zone_id" {
  description = "Internal API's LB Zone ID"
  type        = "string"
}
