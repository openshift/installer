variable "ssh_key" {
  type = "string"
}

variable "vpc_id" {
  type = "string"
}

variable "cl_channel" {
  type = "string"
}

variable "cluster_id" {
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

variable "sg_ids" {
  type        = "list"
  description = "The security group IDs to be applied."
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

variable "root_volume_type" {
  type        = "string"
  description = "The type of volume for the root block device."
}

variable "root_volume_size" {
  type        = "string"
  description = "The size of the volume in gigabytes for the root block device."
}

variable "root_volume_iops" {
  type        = "string"
  default     = "100"
  description = "The amount of provisioned IOPS for the root block device."
}

variable "worker_iam_role" {
  type        = "string"
  default     = ""
  description = "IAM role to use for the instance profiles of worker nodes."
}
