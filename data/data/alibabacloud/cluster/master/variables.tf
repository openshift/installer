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

variable "vswitch_ids" {
  type        = list(string)
  description = "The VSwitch IDs of the master ECS. Example: [vsw-xxx1, vsw-xxx2, vsw-xxx3]"
}

variable "sg_id" {
  type        = string
  description = "The security group ID of the master ECS."
}

variable "slb_ids" {
  type = list(string)
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
  description = "The system disk category of the master ECS.Valid values are cloud_efficiency, cloud_ssd, cloud_essd. Default value is cloud_essd."
  default     = "cloud_essd"
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
  default     = {}
  description = "Tags to be applied to created resources."
}
