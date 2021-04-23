variable "name" {
  type = string
}

variable "instance_count" {
  type = number
}

variable "ignition" {
  type    = string
  default = ""
}

variable "resource_pool" {
  type = string
}

variable "folder" {
  type = string
}

variable "datastore" {
  type = string
}

variable "network" {
  type = string
}

variable "cluster_domain" {
  type = string
}

variable "datacenter" {
  type = string
}

variable "template" {
  type = string
}

variable "guest_id" {
  type = string
}

variable "memory" {
  type = number
}

variable "num_cpus" {
  type = number
}

variable "cores_per_socket" {
  type = number
}

variable "disk_size" {
  type = number
}

variable "tags" {
  type = list(any)
}

variable "cluster_id" {
  type = string
}

variable "thin_disk" {
  type = bool
}

variable "scrub_disk" {
  type = bool
}
