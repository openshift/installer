variable "public_zone_name" {
  description = "The name of the public managed DNS zone"
  type        = string
}

variable "network" {
  description = "URL of the VPC network resource for the cluster"
  type        = string
}

variable "cluster_id" {
  type        = string
  description = "The identifier for the cluster."
}

variable "api_external_lb_ip" {
  description = "External API's LB IP"
  type        = string
}

variable "api_internal_lb_ip" {
  description = "Internal API's LB IP"
  type        = string
}

variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = string
}

variable "public_endpoints" {
  type        = bool
  description = "If the cluster should have externally accessible resources."
}

variable "private_zone_name" {
  description = "The name of the private managed DNS zone"
  type        = string
}

variable "create_private_zone" {
  type        = bool
  description = "Create a private managed zone."
}

variable "create_private_zone_records" {
  type        = bool
  description = "Create records for the private managed zone."
}

variable "create_public_zone_records" {
  type        = bool
  description = "Create records for the public managed zone."
}

variable "public_zone_project" {
  type        = string
  description = "Project where the public managed zone will exist."
}

variable "private_zone_project" {
  type        = string
  description = "Project where the private managed zone will exist."
}