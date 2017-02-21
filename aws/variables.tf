variable "external_vpc_id" {
  type = "string"
}

variable "az_count" {
  type = "string"
}

variable "etcd_domain" {
  type    = "string"
  default = "etcd.cluster."
}
