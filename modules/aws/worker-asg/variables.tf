variable "ssh_key" {
  type = "string"
}

variable "vpc_id" {
  type = "string"
}

variable "cl_channel" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "ec2_type" {
  type = "string"
}

variable "instance_count" {
  type = "string"
}

variable "subnet_ids" {
  type = "list"
}

variable "extra_sg_ids" {
  type = "list"
}

variable "user_data" {
  type        = "string"
  description = "User-data content used to boot the instances"
}