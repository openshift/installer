variable "azure_environment" {
  type        = string
  description = "The target Azure cloud environment for the cluster."
}

variable "azure_region" {
  type        = string
  description = "The target Azure region for the cluster."
}

variable "azure_master_vm_type" {
  type        = string
  description = "Instance type for the master node(s). Example: `Standard_D8s_v3`."
}

variable "azure_master_disk_encryption_set_id" {
  type        = string
  default     = null
  description = "The ID of the Disk Encryption Set which should be used to encrypt OS disk for the master node(s)."
}

variable "azure_master_encryption_at_host_enabled" {
  type        = bool
  description = "Enables encryption at the VM host for the master node(s)."
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
  type = string
  description = "The type of the volume the root block device of master nodes."
}

variable "azure_master_root_volume_size" {
  type = string
  description = "The size of the volume in gigabytes for the root block device of master nodes."
}

variable "azure_control_plane_ultra_ssd_enabled" {
  type = bool
  description = "Determines if the control plane should have UltraSSD Enabled."
}

variable "azure_base_domain_resource_group_name" {
  type = string
  default = ""
  description = "The resource group that contains the dns zone used as base domain for the cluster."
}

variable "azure_image_url" {
  type = string
  description = "The URL of the vm image used for all nodes."
}

variable "azure_arm_endpoint" {
  type = string
  description = "The endpoint for the Azure API. Only used when installing to Azure Stack"
}

variable "azure_bootstrap_ignition_stub" {
  type = string
  description = "The bootstrap ignition stub. Only used when installing to Azure Stack"
}

variable "azure_bootstrap_ignition_url_placeholder" {
  type = string
  description = <<EOF
The placeholder value in the bootstrap ignition to be replaced with the ignition URL.
Only used when installing to Azure Stack
EOF
}

variable "azure_subscription_id" {
  type        = string
  description = "The subscription that should be used to interact with Azure API"
}

variable "azure_client_id" {
  type        = string
  description = "The app ID that should be used to interact with Azure API"
  default     = ""
}

variable "azure_client_secret" {
  type        = string
  description = "The password that should be used to interact with Azure API"
  default     = ""
}

variable "azure_certificate_path" {
  type        = string
  description = "The location of the Azure Service Principal client certificates"
  default     = ""
}

variable "azure_certificate_password" {
  type        = string
  description = "The password for the provided Azure Service Principal client certificates"
  default     = ""
}

variable "azure_tenant_id" {
  type        = string
  description = "The tenant ID that should be used to interact with Azure API"
}

variable "azure_use_msi" {
  type        = bool
  default     = false
  description = "Specifies if we are to use a managed identity for authentication"
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

variable "azure_resource_group_name" {
  type        = string
  description = <<EOF
The name of the resource group for the cluster. If this is set, the cluster is installed to that existing resource group
otherwise a new resource group will be created using cluster id.
EOF
}

variable "azure_network_resource_group_name" {
  type = string
  description = "The name of the network resource group, either existing or to be created."
}

variable "azure_virtual_network" {
  type = string
  description = "The name of the virtual network, either existing or to be created."
}

variable "azure_control_plane_subnet" {
  type = string
  description = "The name of the subnet for the control plane, either existing or to be created."
}

variable "azure_compute_subnet" {
  type = string
  description = "The name of the subnet for worker nodes, either existing or to be created"
}

variable "azure_private" {
  type = bool
  description = "This determines if this is a private cluster or not."
}

variable "azure_outbound_routing_type" {
  type = string
  default = "Loadbalancer"

  description = <<EOF
This determined the routing that will be used for egress to Internet.
When not set, Standard LB will be used for egress to the Internet.
EOF
}

variable "azure_hypervgeneration_version" {
  type        = string
  description = <<EOF
This determines the HyperVGeneration disk type to use for the control plane VMs.
EOF
}

variable "azure_control_plane_vm_networking_type" {
  type = bool
  description = "Whether to enable accelerated networking on control plane nodes."
}

variable "random_storage_account_suffix" {
  type = string
  description = "A random string generated to add a suffix to the storage account and blob"
}

variable "azure_vm_architecture" {
  type = string
  description = "Architecture of the VMs - used when creating images in the image gallery"
}

variable "azure_image_release" {
  type = string
  description = "RHCOS release image version - used when creating the image definition in the gallery"
}

variable "azure_use_marketplace_image" {
  type = bool
  description = "Whether to use a Marketplace image for all nodes"
}

variable "azure_marketplace_image_has_plan" {
  type = bool
  description = "Whether the Marketplace image has a purchase plan"
}

variable "azure_marketplace_image_publisher" {
  type = string
  description = "Publisher of the marketplace image"
  default = ""
}

variable "azure_marketplace_image_offer" {
  type = string
  description = "Offer of the marketplace image"
  default = ""
}

variable "azure_marketplace_image_sku" {
  type = string
  description = "SKU of the marketplace image"
  default = ""
}

variable "azure_marketplace_image_version" {
  type = string
  description = "Version of the marketplace image"
  default = ""
}

variable "azure_master_security_encryption_type" {
  type = string
  default = null

  description = <<EOF
Defines the encryption type when the Virtual Machine is a Confidential VM. Possible values are VMGuestStateOnly and DiskWithVMGuestState.
When set to "VMGuestStateOnly" azure_master_vtpm_enabled should be set to true.
When set to "DiskWithVMGuestState" both azure_master_vtp_enabled and azure_master_secure_boot_enabled should be true.
EOF
}

variable "azure_master_secure_vm_disk_encryption_set_id" {
  type    = string
  default = null

  description = <<EOF
Defines the ID of the Disk Encryption Set which should be used to encrypt this OS Disk when the Virtual Machine is a Confidential VM.
It can only be set when azure_master_security_encryption_type is set to "DiskWithVMGuestState".
EOF
}

variable "azure_master_secure_boot" {
  type = string
  description = "Defines whether the instance should have secure boot enabled."
  default = ""
}

variable "azure_master_virtualized_trusted_platform_module" {
  type = string
  description = "Defines whether the instance should have vTPM enabled."
  default = ""
}

variable "azure_keyvault_resource_group" {
  type = string
  description = "Defines the resource group of the key vault used for storage account encryption."
  default = ""
}

variable "azure_keyvault_name" {
  type = string
  description = "Defines the name of the key vault used for storage account encryption."
  default = ""
}

variable "azure_keyvault_key_name" {
  type = string
  description = "Defines the key in the key vault used for storage account encryption."
  default = ""
}

variable "azure_user_assigned_identity_key" {
  type = string
  description = "Defines the user identity key used for the encryption of storage account."
  default = ""
}

variable "azure_resource_group_metadata_tags" {
  type = map(string)
  description = "Metadata Azure tags to be applied to the cluster resource group."
  default = {}
}
