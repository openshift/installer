variable "control_plane_ips" {
  type    = list(string)
  default = []
}

variable "bootstrap_ip" {
  type    = string
  default = null
}
