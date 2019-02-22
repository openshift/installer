variable "cluster_id" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "dns_server_ip" {
  type    = "string"
  default = ""
}

variable "instance_type" {
  type = "string"
}

variable "ec2_ami" {
  type    = "string"
  default = ""
}

variable "instance_count" {
  type = "string"
}

variable "kubeconfig_content" {
  type    = "string"
  default = ""
}

variable "master_sg_ids" {
  type        = "list"
  description = "The security group IDs to be applied to the master nodes."
}

variable "root_volume_iops" {
  type        = "string"
  description = "The amount of provisioned IOPS for the root block device."
}

variable "root_volume_size" {
  type        = "string"
  description = "The size of the volume in gigabytes for the root block device."
}

variable "root_volume_type" {
  type        = "string"
  description = "The type of volume for the root block device."
}

variable "subnet_ids" {
  type = "list"
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "AWS tags to be applied to created resources."
}

variable "target_group_arns" {
  type        = "list"
  default     = []
  description = "The list of target group ARNs for the load balancer."
}

variable "target_group_arns_length" {
  description = "The length of the 'target_group_arns' variable, to work around https://github.com/hashicorp/terraform/issues/12570."
}

variable "user_data_ign" {
  type = "string"
}
