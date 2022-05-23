variable "resource_pool" {
  type    = map(any)
  default = {}
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

variable "datastore" {
  type    = map(any)
  default = {}
}

variable "datacenter" {
  type    = map(any)
  default = {}
}

variable "template" {
  type    = map(any)
  default = {}
}

variable "tags" {
  type = list(any)
}
