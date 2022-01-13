variable "cluster_id" {
  type = string
}

variable "resource_group_id" {
  type = string
}

variable "vpc_id" {
  type        = string
  description = "The VPC ID of the master ECS."
}

variable "zone_ids" {
  type        = list(string)
  description = "The availability zones in which to create the masters. The length of this list must match instance_count."
}

variable "az_to_vswitch_id" {
  type        = map(string)
  description = "Map from availability zone ID to the ID of the VSwitch in that availability zone"
}

variable "sg_id" {
  type        = string
  description = "The security group ID of the master ECS."
}

variable "slb_ids" {
  type = list(string)
}

variable "slb_group_length" {
  description = "The length of the 'slb_ids' variable, to work around https://github.com/hashicorp/terraform/issues/12570."
}

variable "instance_count" {
  type = string
}

variable "instance_type" {
  type        = string
  description = "The instance type of the master ECS."
}

variable "image_id" {
  type        = string
  description = "The image ID of the master ECS."
}

variable "system_disk_size" {
  type        = number
  description = "The system disk size of the master ECS."
}

variable "system_disk_category" {
  type        = string
  description = "The system disk category of the master ECS. Valid values are cloud_efficiency, cloud_ssd, cloud_essd."
}

variable "role_name" {
  type        = string
  description = "Instance RAM role name. The name is provided and maintained by RAM."
}

variable "user_data_ign" {
  type = string
}

variable "tags" {
  type        = map(string)
  description = "Tags to be applied to created resources."
}

variable "publish_strategy" {
  type        = string
  description = "The cluster publishing strategy, either Internal or External"
}
