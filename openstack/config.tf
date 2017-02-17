variable "flavor_id" {
  type    = "string"
  default = "5cf64088-893b-46b5-9bb1-ee020277635d"
}

variable "image_id" {
  type    = "string"
  default = "3a0c0bac-fa91-4c96-bfcb-ee215ba1cd4d"
}

variable "tectonic_version" {
  type    = "string"
  default = "v1.5.2_coreos.1"
}

variable "controller_count" {
  type    = "string"
  default = "1"
}

variable "worker_count" {
  type    = "string"
  default = "3"
}

variable "etcd_count" {
  type    = "string"
  default = "1"
}

variable "base_domain" {
  type    = "string"
  default = "openstack.dev.coreos.systems"
}

variable "cluster_name" {
  type    = "string"
  default = "demo"
}
