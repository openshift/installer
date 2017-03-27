variable "tectonic_aws_ssh_key" {
  type = "string"
}

variable "tectonic_aws_master_ec2_type" {
  type = "string"
}

variable "tectonic_aws_worker_ec2_type" {
  type = "string"
}

variable "tectonic_aws_etcd_ec2_type" {
  type = "string"
}

variable "tectonic_aws_vpc_cidr_block" {
  type = "string"
  default = "10.0.0.0/16"
}

variable "tectonic_aws_az_count" {
  type = "string"
}

variable "tectonic_aws_external_vpc_id" {
  type = "string"
  default = ""
}

variable "tectonic_aws_external_master_subnet_ids" {
  type    = "list"
  default = [""]
}

variable "tectonic_aws_external_worker_subnet_ids" {
  type    = "list"
  default = [""]
}
