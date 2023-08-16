variable "public_zone_name" {
  description = "The name of the public managed DNS zone"
  type        = string
}

variable "private_zone_name" {
  description = "The name of the private managed DNS zone, if exists"
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

variable "project_id" {
  type        = string
  description = "The target GCP project for the cluster."
}

variable "gcp_extra_labels" {
  type        = map(string)
  description = "GCP labels to be applied to created resources."
  default     = {}
}
