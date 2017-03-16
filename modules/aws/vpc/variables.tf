variable "tectonic_aws_az_count" {
  type = "string"
}

variable "tectonic_aws_external_vpc_id" {
  type = "string"
}

variable "tectonic_aws_vpc_cidr_block" {
  type = "string"
}

variable "tectonic_cluster_name" {
  type = "string"
}

variable "tectonic_aws_external_vpc_master_subnets" {
  type    = "list"
  default = ["a", "b", "c"]
}

variable "tectonic_aws_external_vpc_worker_subnets" {
  type    = "list"
  default = ["a", "b", "c"]
}
