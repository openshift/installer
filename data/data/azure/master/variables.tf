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

variable "vm_image" {
  type        = string
  description = "The resource id of the vm image used for masters."
}

variable "identity" {
  type        = string
  description = "The user assigned identity id for the vm."
}

variable "instance_count" {
  type = string
}

variable "elb_backend_pool_v4_id" {
  type = string
}

variable "elb_backend_pool_v6_id" {
  type = string
}

variable "ilb_backend_pool_v4_id" {
  type = string
}

variable "ilb_backend_pool_v6_id" {
  type = string
}

variable "ignition_master" {
  type    = string
  default = ""
}

variable "kubeconfig_content" {
  type    = string
  default = ""
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

variable "availability_zones" {
  type        = list(string)
  description = "List of the availability zones in which to create the masters. The length of this list must match instance_count."
}

variable "private" {
  type        = bool
  description = "This value determines if this is a private cluster or not."
}

variable "use_ipv4" {
  type        = bool
  description = "This value determines if this is cluster should use IPv4 networking."
}

variable "use_ipv6" {
  type        = bool
  description = "This value determines if this is cluster should use IPv6 networking."
}

variable "emulate_single_stack_ipv6" {
  type        = bool
  description = "This determines whether a dual-stack cluster is configured to emulate single-stack IPv6."
}
