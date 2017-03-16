variable "ssh_key" {
  type = "string"
}

variable "vpc_id" {
  type = "string"
}

variable "tectonic_base_domain" {
  type = "string"
}

variable "tectonic_cluster_name" {
  type = "string"
}

variable "tectonic_cl_channel" {
  type = "string"
}

variable "tectonic_aws_master_ec2_type" {
  type = "string"
}

variable "tectonic_master_count" {
  type = "string"
}

variable "etcd_endpoints" {
  type = "list"
}

variable "master_subnet_ids" {
  type = "list"
}

variable "extra_sg_ids" {
  type = "list"
}

variable "kube_image_url" {
  type = "string"
}

variable "kube_image_tag" {
  type = "string"
}
