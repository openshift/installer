variable "elb_backend_pool_v4_id" {
  type        = string
  default     = null
  description = "The external load balancer bakend pool id. used to attach the bootstrap NIC"
}

variable "ilb_backend_pool_v4_id" {
  type        = string
  default     = null
  description = "The internal load balancer bakend pool id. used to attach the bootstrap NIC"
}

variable "elb_pip_v4" {
  type    = string
  default = null
}

variable "elb_pip_v4_fqdn" {
  type    = string
  default = null
}

variable "ilb_ip_v4_address" {
  type = string
}

variable "virtual_network_id" {
  description = "The ID for Virtual Network that will be linked to the Private DNS zone."
  type        = string
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
  description = "The resource id of the vm image used for bootstrap."
}

variable "availability_set_id" {
  type        = string
  description = "ID of the availability set in which to place VMs"
}

variable "bootstrap_ip" {
  type        = string
  description = "The ip of the bootstrap node. Used for log gathering but not for infrastructure provisioning."
}