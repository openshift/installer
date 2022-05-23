variable "resource_pool" {
  type = map(any)
}

variable "bootstrap_moid" {
  type    = string
  default = ""
}

variable "datastore" {
  type = map(any)
}

variable "datacenter" {
  type = map(any)
}

variable "tags" {
  type = list(any)
}

variable "template" {
  type    = map(any)
  default = {}
}
