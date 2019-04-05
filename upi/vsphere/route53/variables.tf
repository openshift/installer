variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = "string"
}

variable "base_domain" {
  description = "The base domain used for public records."
  type        = "string"
}

variable "bootstrap_count" {
  type = "string"
}

variable "bootstrap_ips" {
  type = "list"
}

variable "control_plane_count" {
  type = "string"
}

variable "control_plane_ips" {
  type = "list"
}

variable "compute_count" {
  type = "string"
}

variable "compute_ips" {
  type = "list"
}
