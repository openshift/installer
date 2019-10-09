variable "name" {
  type = "string"
}

variable "instance_count" {
  type = "string"
}

variable "ignition" {
  type    = "string"
  default = ""
}

variable "cluster_domain" {
  type = "string"
}

variable "machine_cidr" {
  type = "string"
}

variable "ipam" {
  type = "string"
}

variable "ipam_token" {
  type = "string"
}

variable "ip_addresses" {
  type = "list"
}
