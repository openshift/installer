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

variable "extra_tags" {
  description = "Extra AWS tags to be applied to created resources."
  type        = "map"
  default     = {}
}

variable "autoscaling_group_extra_tags" {
  description = "Extra AWS tags to be applied to created autoscaling group resources."
  type        = "list"
  default     = []
}
