variable "elb_backend_pool_v4_id" {
  type        = string
  description = "The external load balancer backend pool id. used to attach the bootstrap NIC"
}

variable "ilb_backend_pool_v4_id" {
  type        = string
  description = "The internal load balancer backend pool id. used to attach the bootstrap NIC"
}

variable "master_subnet_id" {
  type        = string
  description = "The subnet ID for the bootstrap node."
}

variable "nsg_name" {
  type        = string
  description = "The network security group for the subnet."
}

variable "resource_group_name" {
  type        = string
  description = "The resource group name for the deployment."
}

variable "storage_account" {
  type        = any
  description = "the storage account for the cluster. It can be used for boot diagnostics."
}

variable "vm_image" {
  type        = string
  description = "The URI of the vm image to used for bootstrap."
}

variable "availability_set_id" {
  type        = string
  description = "ID of the availability set in which to place VMs"
}
