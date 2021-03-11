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

variable "datacenter" {
  type = string
}

variable "template" {
  type = string
}

variable "guest_id" {
  type = string
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