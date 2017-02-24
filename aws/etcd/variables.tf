variable "tectonic_domain" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "dns_zone" {
  type = "string"
}

variable "node_count" {
  default = "3"
}

variable "vpc_id" {
  type = "string"
}

variable "ssh_key" {
  type = "string"
}

variable "coreos_ami" {
  type = "string"
}

variable "etcd_subnets" {
  type = "list"
}
