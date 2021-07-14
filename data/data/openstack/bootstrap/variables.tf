variable "master_sg_ids" {
  type        = list(string)
  description = "The security group IDs to be applied to the master nodes."
}

variable "private_network_id" {
  type = string
}

variable "nodes_subnet_id" {
  type = string
}

variable "master_port_ids" {
  type = list(string)
}
