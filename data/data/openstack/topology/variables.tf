variable "cidr_block" {
  type = "string"
}

variable "cluster_id" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "control_plane_count" {
  type = "string"
}

variable "external_network" {
  description = "UUID of the external network providing Floating IP addresses."
  type        = "string"
  default     = ""
}

variable "lb_floating_ip" {
  description = "(optional) Existing floating IP address to attach to the load balancer created by the installer."
  type        = "string"
  default     = ""
}

variable "trunk_support" {
  type = "string"
}
