variable "external_vpc_id" {
  type = "string"
}

variable "vpc_cid_block" {
  type    = "string"
  default = "172.31.0.0/16"
}

variable "az_count" {
  type = "string"
}
