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

variable "boot_diag_blob_endpoint" {
  type        = string
  description = "the blob endpoint where machines should store their boot diagnostics."
}

variable "ssh_nat_rule_id" {
  type        = string
  description = "ssh nat rule to make the bootstrap node reachable"
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "tags to be applied to created resources."
}

variable "private_dns_zone_id" {
  type        = string
  description = "This is to create explicit dependency on private zone to exist before VMs are created in the vnet. https://github.com/MicrosoftDocs/azure-docs/issues/13728"
}
