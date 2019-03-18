variable "cluster_domain" {
  description = "The domain for the cluster that all DNS records must belong"
  type        = "string"
}

variable "bootstrap_ip" {
  type = "string"
}

variable "control_plane_ips" {
  type = "list"
}

variable "base_domain" {
  description = "The base domain used for public records."
  type        = "string"
}
