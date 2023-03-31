variable "resource_group_name" {
  type        = string
  description = "The resource group name for the deployment."
}

variable "storage_account_name" {
  type        = string
  description = "the name of the storage account for the cluster. It can be used for boot diagnostics."
}

variable "rhcos_image_url" {
  type        = string
  description = "The external load balancer bakend pool id. used to attach the bootstrap NIC"
}