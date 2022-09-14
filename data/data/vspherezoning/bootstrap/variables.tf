variable "resource_pool" {
  type    = list(any)
  default = []
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
  type    = list(any)
  default = []
}

variable "datacenter" {
  type    = list(any)
  default = []
}

variable "template" {
  type    = list(any)
  default = []
}

variable "tags" {
  type = list(any)
}
