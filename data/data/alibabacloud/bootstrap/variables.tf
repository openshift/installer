variable "cluster_id" {
  type = string
}

variable "resource_group_id" {
  type = string
}

variable "ignition" {
  type        = string
  description = "The path of the bootstrap ignition file."
}

variable "ignition_bucket" {
  type        = string
  description = "The name of the new OSS bucket."
}

variable "ignition_stub" {
  type        = string
  description = <<EOF
The stub Ignition config that should be used to boot the bootstrap instance. This already points to the presigned URL for the OSS bucket
specified in ignition_bucket.
EOF
}

variable "vpc_id" {
  type = string
  description = "The VPC id of the bootstrap ECS."
}

variable "vswitch_id" {
  type = string
  description = "The VSwitch id of the bootstrap ECS."
}

variable "slb_id" {
  type = string
  description = "The load balancer of the bootstrap ECS."
}

variable "instance_type" {
  type = string
  description = "The instance type of the bootstrap ECS."
}

variable "image_id" {
  type = string
  description = "The image id of the bootstrap ECS."
}

variable "system_disk_size" {
  type = number
  description = "The system disk size of the bootstrap ECS."
}

variable "system_disk_category" {
  type = string
  description = "The system disk category of the bootstrap ECS.Valid values are cloud_efficiency, cloud_ssd, cloud_essd. Default value is cloud_essd."
  default = "cloud_essd"
}

variable "key_name" {
  type = string
  description = "The name of key pair that can login bootstrap ECS instance successfully without password."
}

variable "tags" {
  type = map(string)
  default = {}
  description = "Tags to be applied to created resources."
}
