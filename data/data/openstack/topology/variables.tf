variable "cidr_block" {
  type = string
}

variable "cluster_id" {
  type = string
}

variable "external_network" {
  description = "Name of the external network providing Floating IP addresses."
  type        = string
  default     = ""
}

variable "external_network_id" {
  description = "UUID of the external network providing Floating IP addresses."
  type        = string
  default     = ""
}

variable "lb_floating_ip" {
  description = "(optional) Existing floating IP address to attach to the load balancer created by the installer."
  type        = string
  default     = ""
}

variable "masters_count" {
  type = string
}

variable "api_int_ip" {
  type = string
}

variable "node_dns_ip" {
  type = string
}

variable "trunk_support" {
  type = string
}

variable "octavia_support" {
  type = "string"
}
