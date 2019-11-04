variable "nsx_manager" {
  type        = "string"
  description = "This is the NSX manager for the environment."
}

variable "nsx_user" {
  type        = "string"
  description = "NSX manager user for the environment."
}

variable "nsx_password" {
  type        = "string"
  description = "NSX manager password"
}

variable "nsx_tag_scope" {
  default = "project"
}

variable "nsx_tag" {
  default = "terraform-demo"
}

variable "t1_router" {
  type = "string"
}

variable "logical_switch" {
  type = "string"
}
variable "nsx_edge_cluster" {
  type = "string"
}

variable "datacenter_id" {
  type = "string"
}

variable "dhcp_server_ip" {
  type = "string"
}

variable "gateway_ip" {
  type = "string"
}

variable "ip_block_cidr" {
  type = "string"
}

variable "ip_pool_start" {
  type = "string"
}

variable "ip_pool_end" {
  type = "string"
}

variable "dns_nameservers" {
  type = "string"
}

variable "logical_switch_ip_address" {
  type = "string"
}

variable "ip_pool_cidr" {
  type = "string"
}

variable "base_domain" {
  type = "string"
}

variable "transport_zone_id" {
  type = "string"
}