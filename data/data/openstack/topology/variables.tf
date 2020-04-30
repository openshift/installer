variable "cidr_block" {
  type = string
}

variable "cluster_id" {
  type = string
}

variable "cluster_domain" {
  type        = string
  description = "The domain name of the cluster. All DNS records must be under this domain."
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

variable "ingress_ip" {
  type = string
}

variable "external_dns" {
  type = list(string)
}

variable "trunk_support" {
  type = string
}

variable "octavia_support" {
  type = string
}

variable "machines_subnet_id" {
  type    = string
  default = ""
}

variable "machines_network_id" {
  type    = string
  default = ""
}
