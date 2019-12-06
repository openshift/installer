variable "vm_size" {
  type        = string
  description = "The SKU ID for the bootstrap node."
}

variable "vm_image" {
  type        = string
  description = "The resource id of the vm image used for bootstrap."
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

variable "identity" {
  type        = string
  description = "The user assigned identity id for the vm."
}

variable "ignition" {
  type        = string
  description = "The content of the bootstrap ignition file."
}

variable "subnet_id" {
  type        = string
  description = "The subnet ID for the bootstrap node."
}

variable "elb_backend_pool_id" {
  type        = string
  description = "The external load balancer bakend pool id. used to attach the bootstrap NIC"
}

variable "ilb_backend_pool_id" {
  type        = string
  description = "The internal load balancer bakend pool id. used to attach the bootstrap NIC"
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
