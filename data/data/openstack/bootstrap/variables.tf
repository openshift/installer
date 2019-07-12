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

variable "bootstrap_port_id" {
  type        = string
  description = "The subnet ID for the bootstrap node."
}
