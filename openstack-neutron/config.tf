variable "flavor_id" {
  type    = "string"
  default = "5cf64088-893b-46b5-9bb1-ee020277635d"
}

variable "image_id" {
  type    = "string"
  default = "3acad946-7dd9-487d-b76f-75c79b8d550b"
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

variable "public_network_name" {
  type    = "string"
  default = "public"
}

variable "external_gateway_id" {
  type    = "string"
  default = "6d6357ac-0f70-4afa-8bd7-c274cc4ea235"
}
