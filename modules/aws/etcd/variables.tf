variable "tectonic_base_domain" {
  type = "string"
}

variable "tectonic_cluster_name" {
  type = "string"
}

variable "tectonic_cl_channel" {
  type = "string"
}

variable "dns_zone" {
  type = "string"
}

variable "az_count" {
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

variable "etcd_subnets" {
  type = "list"
}

variable "external_endpoints" {
  type = "list"
}

variable "etcd_version" {
  type = "string"
}
