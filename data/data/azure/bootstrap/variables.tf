variable "vm_size" {
  type        = "string"
  description = "The SKU ID for the bootstrap node."
}

variable "nsg_id" {
  type  = "string"
  description = "The NSG attached to the subnet"
}

variable "region" {
  type        = "string"
  description = "The region for the deployment."
}

variable "resource_group_name" {
  type        = "string"
  description = "The resource group name for the deployment."
}

variable "cluster_id" {
  type        = "string"
  description = "The identifier for the cluster."
}

variable "ignition" {
  type        = "string"
  description = "The content of the bootstrap ignition file."
}

variable "subnet_id" {
  type        = "string"
  description = "The subnet ID for the bootstrap node."
}

variable "elb_backend_pool_id" {
  type        = "string"
  description ="The external load balancer bakend pool id. used to attach the bootstrap NIC"
}

variable "ilb_backend_pool_id" {
  type        = "string"
  description ="The internal load balancer bakend pool id. used to attach the bootstrap NIC"
}

variable "external_lb_id" {
  type        = "string"
  description ="The external load balancer id"
}

variable "boot_diag_blob_endpoint" {
  type = "string"
  description = "the blob endpoint where machines should store their boot diagnostics."
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "AWS tags to be applied to created resources."
}

variable "target_group_arns" {
  type        = "list"
  default     = []
  description = "The list of target group ARNs for the load balancer."
}

variable "volume_iops" {
  type        = "string"
  default     = "100"
  description = "The amount of IOPS to provision for the disk."
}

variable "volume_size" {
  type        = "string"
  default     = "30"
  description = "The volume size (in gibibytes) for the bootstrap node's root volume."
}
