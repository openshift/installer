variable "azure_config_version" {
  description = <<EOF
(internal) This declares the version of the Azure configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF


  default = "0.1"
}

variable "azure_region" {
  type = string
  description = "The target Azure region for the cluster."
}

variable "azure_bootstrap_vm_type" {
  type = string
  description = "Instance type for the bootstrap node. Example: `Standard_DS4_v3`."
}

variable "azure_master_vm_type" {
  type = string
  description = "Instance type for the master node(s). Example: `Standard_DS4_v3`."
}

variable "azure_extra_tags" {
  type = map(string)

  description = <<EOF
(optional) Extra Azure tags to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF


default = {}
}

variable "azure_master_root_volume_type" {
  type        = string
  description = "The type of the volume the root block device of master nodes."
}

variable "azure_master_root_volume_size" {
  type        = string
  description = "The size of the volume in gigabytes for the root block device of master nodes."
}

variable "azure_base_domain_resource_group_name" {
  type        = string
  description = "The resource group that contains the dns zone used as base domain for the cluster."
}

variable "azure_image_url" {
  type        = string
  description = "The URL of the vm image used for all nodes."
}

variable "azure_subscription_id" {
  type        = string
  description = "The subscription that should be used to interact with Azure API"
}

variable "azure_client_id" {
  type        = string
  description = "The app ID that should be used to interact with Azure API"
}

variable "azure_client_secret" {
  type        = string
  description = "The password that should be used to interact with Azure API"
}

variable "azure_tenant_id" {
  type        = string
  description = "The tenant ID that should be used to interact with Azure API"
}

variable "azure_master_availability_zones" {
  type        = list(string)
  description = "The availability zones in which to create the masters. The length of this list must match master_count."
}

variable "azure_preexisting_network" {
  type        = bool
  default     = false
  description = "Specifies whether an existing network should be used or a new one created for installation."
}

variable "azure_network_resource_group_name" {
  type        = string
  description = "The name of the network resource group, either existing or to be created."
}

variable "azure_virtual_network" {
  type        = string
  description = "The name of the virtual network, either existing or to be created."
}

variable "azure_control_plane_subnet" {
  type        = string
  description = "The name of the subnet for the control plane, either existing or to be created."
}

variable "azure_compute_subnet" {
  type        = string
  description = "The name of the subnet for worker nodes, either existing or to be created"
}

variable "azure_private" {
  type        = bool
  description = "This determines if this is a private cluster or not."
}

variable "azure_emulate_single_stack_ipv6" {
  type        = bool
  description = "This determines whether a dual-stack cluster is configured to emulate single-stack IPv6."
}
