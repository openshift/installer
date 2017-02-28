variable "flavor_id" {
  type    = "string"
  default = "bbcb7eb5-5c8d-498f-9d7e-307c575d3566"
}

variable "image_id" {
  type    = "string"
  default = "acdcd535-5408-40f3-8e88-ad8ebb6507e6"
}

variable "tectonic_version" {
  type    = "string"
  default = "v1.5.2_coreos.1"
}

variable "master_count" {
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
