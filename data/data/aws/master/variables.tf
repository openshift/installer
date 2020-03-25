variable "availability_zones" {
  type        = list(string)
  description = "List of the availability zones in which to create the masters. The length of this list must match instance_count."
}

variable "az_to_subnet_id" {
  type        = map(string)
  description = "Map from availability zone name to the ID of the subnet in that availability zone"
}

variable "cluster_id" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "ec2_ami" {
  type    = string
  default = ""
}

variable "instance_count" {
  type = string
}

variable "kubeconfig_content" {
  type    = string
  default = ""
}

variable "master_sg_ids" {
  type        = list(string)
  description = "The security group IDs to be applied to the master nodes."
}

variable "root_volume_iops" {
  type        = string
  description = "The amount of provisioned IOPS for the root block device."
}

variable "root_volume_size" {
  type        = string
  description = "The size of the volume in gigabytes for the root block device."
}

variable "root_volume_type" {
  type        = string
  description = "The type of volume for the root block device."
}

variable "root_volume_encrypted" {
  type        = bool
  description = "Whether the root block device should be encrypted."
}

variable "root_volume_kms_key_id" {
  type        = string
  description = "The KMS key id that should be used tpo encrypt the root block device."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "AWS tags to be applied to created resources."
}

variable "target_group_arns" {
  type        = list(string)
  default     = []
  description = "The list of target group ARNs for the load balancer."
}

variable "target_group_arns_length" {
  description = "The length of the 'target_group_arns' variable, to work around https://github.com/hashicorp/terraform/issues/12570."
}

variable "user_data_ign" {
  type = string
}

variable "publish_strategy" {
  type        = string
  description = <<EOF
The publishing strategy for endpoints like load balancers.

Because of the issue https://github.com/hashicorp/terraform/issues/12570, the consumers cannot use a dynamic list for count
and therefore are force to implicitly assume that the list is of aws_lb_target_group_arns_length - 1, in case there is no api_external. And that's where this variable
helps to decide if the target_group_arns is of length (target_group_arns_length) or (target_group_arns_length - 1)
EOF
}
