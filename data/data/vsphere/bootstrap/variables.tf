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

variable "tags" {
  type = list(any)
}
