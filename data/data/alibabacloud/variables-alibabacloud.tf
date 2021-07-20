variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}

variable "region_id" {
  type        = string
  description = "The target Alibaba Cloud region for the cluster."
}

variable "zone_ids" {
  type        = list(string)
  description = "The availability zones in which to create the masters and workers."
}

variable "resource_group_id" {
  type = string
}

variable "instance_type" {
  type        = string
  description = "The instance type of the master ECS."
}

variable "image_id" {
  type        = string
  description = "The image id of the master ECS."
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

variable "key_name" {
  type        = string
  description = "The name of key pair that can login ECS instance successfully without password."
}

variable "resource_tags" {
  type = map(string)

  description = <<EOF
(optional) Extra tags to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF
  default = {}
}

variable "ignition_bucket" {
  type = string
  description = "The name of the new OSS bucket."
}