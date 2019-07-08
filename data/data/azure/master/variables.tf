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

variable "external_lb_id" {
  type = string
}

variable "elb_backend_pool_id" {
  type = string
}

variable "ilb_backend_pool_id" {
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

variable "os_volume_size" {
  type        = string
  description = "The size of the volume in gigabytes for the root block device."
  default     = "100"
}

variable "tags" {
  type        = map(string)
  default     = {}
  description = "tags to be applied to created resources."
}

variable "boot_diag_blob_endpoint" {
  type        = string
  description = "the blob endpoint where machines should store their boot diagnostics."
}

variable "ignition" {
  type = string
}

variable "master_subnet_cidr" {
  type        = string
  description = "the master subnet cidr"
}

variable "private_dns_zone_id" {
  type        = string
  description = "This is to create explicit dependency on private zone to exist before VMs are created in the vnet. https://github.com/MicrosoftDocs/azure-docs/issues/13728"
}

variable "availability_zones" {
  type        = list(string)
  description = "List of the availability zones in which to create the masters. The length of this list must match instance_count."
}
