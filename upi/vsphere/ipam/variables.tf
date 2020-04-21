variable "hostnames" {
  type = set(string)
}

variable "machine_cidr" {
  type = string
}

variable "ipam" {
  type = string
}

variable "ipam_token" {
  type = string
}

variable "static_ip_addresses" {
  type    = list(string)
  default = []
}
