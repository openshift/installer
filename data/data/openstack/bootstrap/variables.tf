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

variable "cluster_domain" {
  type        = string
  description = "The domain name of the cluster. All DNS records must be under this domain."
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

variable "service_vm_fixed_ip" {
  type = string
}

variable "master_vm_fixed_ip" {
  type        = "string"
  description = "Fixed IP for a master node to provide DNS to bootstrap during clustering."
}
