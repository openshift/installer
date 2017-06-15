variable "tectonic_azure_config_version" {
  description = <<EOF
(internal) This declares the version of the Azure configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF

  default = "1.0"
}

variable "tectonic_azure_ssh_network_internal" {
  type = "string"

  description = <<EOF
Network (internal) to allow SSH access from. Maps to `source_address_prefix` in Azure.
Defaults to `VirtualNetwork`. Should be internal to Azure environment.
Allowed values: [network CIDR (i.e., 10.0.0.0/16) | `VirtualNetwork` | `Internet` | `*` ]
EOF

  default = "VirtualNetwork"
}

variable "tectonic_azure_ssh_network_external" {
  type = "string"

  description = <<EOF
(optional) Network (external) to allow SSH access from. Maps to `source_address_prefix` in Azure.
Defaults to `*`. Can be external to Azure environment.
Allowed values: [network CIDR (i.e., 10.0.0.0/16) | `VirtualNetwork` | `Internet` | `*` ]
EOF

  default = "*"
}

variable "tectonic_azure_dns_resource_group" {
  type    = "string"
  default = "tectonic-dns-group"
}

// The image ID as given in `azure image list`.
// Specifies the OS image of the VM.
variable "tectonic_azure_image_reference" {
  type = "map"

  description = <<EOF
(optional) Specifies an image map with the following keys: `publisher`, `offer`, `sku`, `version`
EOF

  default = {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "latest"
  }
}

variable "tectonic_azure_location" {
  type = "string"
}

// Name of an Azure ssh key to use
// joe-sfo
variable "tectonic_azure_ssh_key" {
  type = "string"
}

variable "tectonic_azure_master_vm_size" {
  type        = "string"
  description = "Instance size for the master node(s). Example: Standard_DS2_v2."
  default     = "Standard_DS2_v2"
}

variable "tectonic_azure_worker_vm_size" {
  type        = "string"
  description = "Instance size for the worker node(s). Example: Standard_DS2_v2."
  default     = "Standard_DS2_v2"
}

variable "tectonic_azure_etcd_vm_size" {
  type        = "string"
  description = "Instance size for the etcd node(s). Example: Standard_DS2_v2."
  default     = "Standard_DS2_v2"
}

variable "tectonic_azure_master_storage_account_type" {
  type        = "string"
  description = "Storage account type for the master node(s). Example: Premium_LRS."
  default     = "Premium_LRS"
}

variable "tectonic_azure_worker_storage_account_type" {
  type        = "string"
  description = "Storage account type for the worker node(s). Example: Premium_LRS."
  default     = "Premium_LRS"
}

variable "tectonic_azure_etcd_storage_account_type" {
  type        = "string"
  description = "Storage account type for the etcd node(s). Example: Premium_LRS."
  default     = "Premium_LRS"
}

variable "tectonic_azure_vnet_cidr_block" {
  type        = "string"
  default     = "10.0.0.0/16"
  description = "Block of IP addresses used by the Resource Group. This should not overlap with any other networks, such as a private datacenter connected via ExpressRoute."
}

variable "tectonic_azure_external_vnet_id" {
  type        = "string"
  description = "ID of an existing Virtual Network to launch nodes into. Example: VNet1. Leave blank to create a new Virtual Network."
  default     = ""
}

variable "tectonic_azure_external_rsg_name" {
  type        = "string"
  default     = ""
  description = "Pre-existing resource group to use as parent for cluster resources."
}

variable "tectonic_azure_external_vnet_name" {
  type        = "string"
  default     = ""
  description = "Pre-existing virtual network to create cluster into."
}

variable "tectonic_azure_create_dns_zone" {
  description = "If set to true, create an Azure DNS zone"
  default     = true
}

variable "tectonic_azure_use_custom_fqdn" {
  description = "(optional) If set to true, assemble the FQDN from the configuration. Otherwise, use the FQDN set up by Azure."
  default     = false
}

variable "tectonic_azure_external_master_subnet_id" {
  type = "string"

  description = <<EOF
(optional) Subnet ID within an existing VNet to deploy master nodes into.
Required to use an existing VNet.

Example: the subnet ID starts with `"/subscriptions/{subscriptionId}"` or `"/providers/{resourceProviderNamespace}"`.
EOF

  default = ""
}

variable "tectonic_azure_external_worker_subnet_id" {
  type = "string"

  description = <<EOF
(optional) Subnet ID within an existing VNet to deploy worker nodes into.
Required to use an existing VNet.

Example: the subnet ID starts with `"/subscriptions/{subscriptionId}"` or `"/providers/{resourceProviderNamespace}"`.
EOF

  default = ""
}

variable "tectonic_azure_external_resource_group" {
  type = "string"

  description = <<EOF
(optional) The name of the resource group of the external Network Security
Group used. This is required if specifying `tectonic_external_nsg_etcd`,
`tectonic_external_nsg_master`, and/or `tectonic_external_nsg_worker`,
EOF

  default = ""
}

variable "tectonic_azure_external_nsg_etcd" {
  type = "string"

  description = <<EOF
(optional) The name of the external Network Security Group used for etcd. This
depends on `tectonic_azure_external_resource_group` to also be specified.
EOF

  default = ""
}

variable "tectonic_azure_external_nsg_api" {
  type = "string"

  description = <<EOF
(optional) The name of the external Network Security Group used for the
Tectonic Console and Kubernetes API Server. This depends on
`tectonic_azure_external_resource_group` to also be specified.
EOF

  default = ""
}

variable "tectonic_azure_external_nsg_master" {
  type = "string"

  description = <<EOF
(optional) The name of the external Network Security Group used for masters. This
depends on `tectonic_azure_external_resource_group` to also be specified.
EOF

  default = ""
}

variable "tectonic_azure_external_nsg_worker" {
  type = "string"

  description = <<EOF
(optional) The name of the external Network Security Group used for workers. This
depends on `tectonic_azure_external_resource_group` to also be specified.
EOF

  default = ""
}
