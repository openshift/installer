variable "ssh_key" {
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

variable "master_sg_ids" {
  type        = "list"
  description = "The security group IDs to be applied to the master nodes."
}

variable "api_sg_ids" {
  type        = "list"
  description = "The security group IDs to be applied to the public facing ELB."
}

variable "console_sg_ids" {
  type        = "list"
  description = "The security group IDs to be applied to the console ELB."
}

variable "base_domain" {
  type        = "string"
  description = "Domain on which the ELB records will be created"
}

variable "internal_zone_id" {
  type        = "string"
  description = "ID of the internal facing Route53 Hosted Zone on which the ELB records will be created"
}

variable "external_zone_id" {
  type        = "string"
  description = "ID of the public facing Route53 Hosted Zone on which the ELB records will be created"
}

variable "user_data" {
  type        = "string"
  description = "User-data content used to boot the instances"
}

variable "public_vpc" {
  description = "If set to true, public facing ingress resources are created."
  default     = true
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

variable "custom_dns_name" {
  type        = "string"
  default     = ""
  description = "DNS prefix used to construct the console and API server endpoints."
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

variable "master_iam_role" {
  type        = "string"
  default     = ""
  description = "IAM role to use for the instance profiles of master nodes."
}
