variable "region" {
  type        = string
  description = "The region for the deployment."
}

variable "resource_group_name" {
  type        = string
  description = "The resource group name for the deployment."
}

variable "cluster_id" {
  type = string
}

variable "vm_size" {
  type = string
}

variable "vm_image_uri" {
  type        = string
  description = "The URI of the vm image used for masters."
}

variable "instance_count" {
  type = string
}

variable "elb_backend_pool_v4_id" {
  type = string
}

variable "ilb_backend_pool_v4_id" {
  type = string
}

variable "subnet_id" {
  type        = string
  description = "The subnet to attach the masters to."
}

variable "os_volume_type" {
  type        = string
  description = "The type of the volume for the root block device."
}

variable "os_volume_size" {
  type        = string
  description = "The size of the volume in gigabytes for the root block device."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "tags to be applied to created resources."
}

variable "storage_account" {
  type        = any
  description = "the storage account for the cluster. It can be used for boot diagnostics."
}

variable "ignition" {
  type = string
}

variable "private" {
  type        = bool
  description = "This value determines if this is a private cluster or not."
}

variable "availability_set_id" {
  type        = string
  description = "ID of the availability set in which to place VMs"
}