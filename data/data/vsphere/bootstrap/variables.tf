variable "resource_pool" {
  type = string
}

variable "bootstrap_moid" {
  type    = string
  default = ""
}

variable "control_plane_ips" {
  type    = list(string)
  default = []
}

variable "control_plane_moids" {
  type    = list(string)
  default = []
}

variable "folder" {
  type = string
}

variable "datastore" {
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
  type = list
}

variable "thin_disk" {
  type = bool
}

variable "scrub_disk" {
  type = bool
}
