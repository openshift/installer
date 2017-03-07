variable "flavor_id" {
  type    = "string"
  default = "5cf64088-893b-46b5-9bb1-ee020277635d"
}

variable "image_id" {
  type    = "string"
  default = "acdcd535-5408-40f3-8e88-ad8ebb6507e6"
}

variable "kube_version" {
  type = "string"
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

variable "external_gateway_id" {
  type    = "string"
  default = "6d6357ac-0f70-4afa-8bd7-c274cc4ea235"
}
