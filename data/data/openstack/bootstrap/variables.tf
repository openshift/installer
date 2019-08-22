variable "image_name" {
  type        = string
  description = "The name of the Glance image for the bootstrap node."
}

variable "swift_container" {
  type        = string
  description = "The Swift container name for bootstrap ignition file."
}

variable "cluster_id" {
  type        = string
  description = "The identifier for the cluster."
}

variable "ignition" {
  type        = string
  description = "The content of the bootstrap ignition file."
}

variable "flavor_name" {
  type        = string
  description = "The Nova flavor for the bootstrap node."
}

variable "api_int_ip" {
  type = string
}

variable "node_dns_ip" {
  type = string
}

variable "external_network" {
  type = string
}

variable "private_network_id" {
  type = string
}

variable "master_sg_id" {
  type = string
}

variable "nodes_subnet_id" {
  type = string
}

variable "cluster_domain" {
  type = string
}
