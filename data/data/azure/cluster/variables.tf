variable "elb_backend_pool_v4_id" {
  type        = string
  default     = null
  description = "The external load balancer bakend pool id. used to attach the bootstrap NIC"
}

variable "elb_backend_pool_v6_id" {
  type        = string
  default     = null
  description = "The external load balancer bakend pool id for ipv6. used to attach the bootstrap NIC"
}

variable "ilb_backend_pool_v4_id" {
  type        = string
  default     = null
  description = "The internal load balancer bakend pool id. used to attach the bootstrap NIC"
}

variable "ilb_backend_pool_v6_id" {
  type        = string
  default     = null
  description = "The internal load balancer bakend pool id for ipv6. used to attach the bootstrap NIC"
}

variable "public_lb_pip_v4_fqdn" {
  type    = string
  default = null
}

variable "public_lb_pip_v6_fqdn" {
  type    = string
  default = null
}

variable "internal_lb_ip_v4_address" {
  type    = string
  default = null
}

variable "internal_lb_ip_v6_address" {
  type    = string
  default = null
}

variable "virtual_network_id" {
  description = "The ID for Virtual Network that will be linked to the Private DNS zone."
  type        = string
}

variable "master_subnet_id" {
  type        = string
  description = "The subnet ID for the bootstrap node."
}

variable "nsg_name" {
  type        = string
  description = "The network security group for the subnet."
}

variable "resource_group_name" {
  type        = string
  description = "The resource group name for the deployment."
}

variable "vm_image" {
  type        = string
  description = "The resource id of the vm image used for bootstrap."
}

variable "identity" {
  type        = string
  description = "The user assigned identity id for the vm."
}

variable "outbound_type" {
  type    = string
  default = "Loadbalancer"

  description = <<EOF
This determined the routing type that will be used for egress to
Internet.
When false, Standard LB will be used for egress to the Internet.

This is required because terraform cannot calculate counts during plan phase
completely and therefore the `vnet/public-lb.tf`
conditional need to be recreated. See
https://github.com/hashicorp/terraform/issues/12570
EOF
}
