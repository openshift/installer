variable "ali_access_key" {
  type = string
}

variable "ali_secret_key" {
  type = string
}

variable "ali_region_id" {
  type        = string
  description = "The target Alibaba Cloud region for the cluster."
}

variable "ali_zone_ids" {
  type        = list(string)
  description = "The availability zones in which to create the masters."
}

variable "ali_nat_gateway_zone_id" {
  type        = string
  description = "The availability zone in which to create the NAT gateway."
}

variable "ali_resource_group_id" {
  type = string
}

variable "ali_bootstrap_instance_type" {
  type        = string
  description = "The instance type the bootstrap ECS."
}

variable "ali_master_instance_type" {
  type        = string
  description = "The instance type of the master ECS."
}

variable "ali_image_id" {
  type        = string
  description = "The image ID of the master ECS."
}

variable "ali_system_disk_size" {
  type        = number
  description = "The system disk size of the master ECS."
}

variable "ali_system_disk_category" {
  type        = string
  description = "The system disk category of the master ECS. Valid values are cloud_efficiency, cloud_ssd, cloud_essd."
}

variable "ali_extra_tags" {
  type = map(string)

  description = <<EOF
(optional) Extra tags to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF
}

variable "ali_ignition_bucket" {
  type = string
  description = "The OSS bucket where the ignition configuration is stored."
}

variable "ali_bootstrap_stub_ignition" {
  type = string
  description = <<EOF
The stub Ignition configuration used to boot the bootstrap ECS instance. This already points to the presigned URL for the OSS bucket
specified in ‘ali_ignition_bucket’.
EOF
}
