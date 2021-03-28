variable "vm_size" {
  type        = string
  description = "The SKU ID for the bootstrap node."
}

variable "vm_image_uri" {
  type        = string
  description = "The URI of the vm image to used for bootstrap."
}

variable "region" {
  type        = string
  description = "The region for the deployment."
}

variable "resource_group_name" {
  type        = string
  description = "The resource group name for the deployment."
}

variable "cluster_id" {
  type        = string
  description = "The identifier for the cluster."
}

variable "ignition" {
  type        = string
  description = "The content of the bootstrap ignition file."
}

variable "subnet_id" {
  type        = string
  description = "The subnet ID for the bootstrap node."
}

variable "elb_backend_pool_v4_id" {
  type        = string
  description = "The external load balancer backend pool id. used to attach the bootstrap NIC"
}

variable "ilb_backend_pool_v4_id" {
  type        = string
  description = "The internal load balancer backend pool id. used to attach the bootstrap NIC"
}

variable "storage_account" {
  type        = any
  description = "the storage account for the cluster. It can be used for boot diagnostics."
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "tags to be applied to created resources."
}

variable "nsg_name" {
  type        = string
  description = "The network security group for the subnet."
}

variable "private" {
  type        = bool
  description = "This value determines if this is a private cluster or not."
}

variable "availability_set_id" {
  type        = string
  description = "ID of the availability set in which to place VMs"
}
