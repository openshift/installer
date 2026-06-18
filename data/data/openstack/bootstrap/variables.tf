variable "openstack_bootstrap_flavor_name" {
  type        = string
  description = "Instance size for the bootstrap node. If not set, the master flavor is used. Example: `m1.medium`."
}

variable "master_sg_ids" {
  type        = list(string)
  description = "The security group IDs to be applied to the master nodes."
}

variable "private_network_id" {
  type = string
}

variable "nodes_default_port" {
  type = object({
    network_id = string
    fixed_ips = list(object({
      subnet_id  = string
      ip_address = string
    }))
  })
}

variable "master_port_ids" {
  type = list(string)
}
