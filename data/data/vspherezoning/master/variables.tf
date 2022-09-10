variable "resource_pool" {
  type = list(any)
}

variable "bootstrap_moid" {
  type    = string
  default = ""
}

variable "datastore" {
  type = list(any)
}

variable "datacenter" {
  type = list(any)
}

variable "tags" {
  type = list(any)
}

variable "template" {
  type    = list(any)
  default = []
}
