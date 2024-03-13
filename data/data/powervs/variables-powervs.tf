################################################################
# Configure the IBM Cloud provider
################################################################
variable "powervs_api_key" {
  type        = string
  description = "IBM Cloud API key associated with user's identity"
  default     = "<key>"
}

variable "powervs_vpc_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = "eu-gb"
}

variable "powervs_vpc_zone" {
  type        = string
  description = "The IBM Cloud zone associated with the VPC region you're using"
}

variable "powervs_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = "lon"
}

variable "powervs_zone" {
  type        = string
  description = "The IBM Cloud zone associated with the region you're using"
}

variable "powervs_resource_group" {
  type        = string
  description = "The cloud instance resource group"
}

variable "powervs_publish_strategy" {
  type        = string
  description = "The cluster publishing strategy, either Internal or External"
  default     = "External"
}

variable "powervs_service_instance_name" {
  type        = string
  description = "Optionally, the service instance name of an existing object before cluster creation"
}

variable "powervs_wait_for_workspace" {
  type        = string
  description = "The seconds wait for the Power VS Workspace after creation, default is 3m"
  default     = "3m"
}

################################################################
# Configure storage
################################################################
variable "powervs_cos_region" {
  type        = string
  description = "The region where your COS instance is located in"
  default     = "eu-gb"
}

variable "powervs_cos_instance_location" {
  type        = string
  description = "The location of your COS instance"
  default     = "global"
}

variable "powervs_cos_storage_class" {
  type        = string
  description = "The plan used for your COS instance"
  default     = "smart"
}

################################################################
# Configure networking
################################################################
variable "powervs_wait_for_vpc" {
  type        = string
  description = "The seconds wait for VPC creation, default is 60s"
  default     = "60s"
}

variable "powervs_vpc_name" {
  type        = string
  description = "The name of a pre-created IBM Cloud VPC. Must be in $powervs_vpc_region"
  default     = ""
}

variable "powervs_vpc_permitted" {
  type        = bool
  description = "Specifies whether an existing VPC is already a Permitted Network for DNS Instance, for Private clusters."
  default     = false
}

variable "powervs_vpc_gateway_attached" {
  type        = bool
  description = "Specifies whether an existing gateway is already attached to an existing VPC."
  default     = false
}

variable "powervs_tg_connection_vpc_id" {
  type        = string
  description = "ID of a VPC connection to the transit gateway specified in attached_transit_gateway, if any."
}

variable "powervs_vpc_gateway_name" {
  type        = string
  description = "The name of a pre-created VPC gateway. Must be in $powervs_vpc_region"
  default     = ""
}

variable "powervs_vpc_subnet_name" {
  type        = string
  description = "The name of a pre-created IBM Cloud Subnet. Must be in $powervs_vpc_region"
  default     = ""
}

variable "powervs_enable_snat" {
  type        = bool
  description = "Boolean indicating if SNAT should be enabled or disabled."
  default     = true
}

variable "powervs_transit_gateway_enabled" {
  type        = bool
  description = "Boolean indicating if Transit Gateways should be used."
  default     = false
}

variable "powervs_attached_transit_gateway" {
  type        = string
  description = "ID of already attached Transit Gateways."
  default     = ""
}

################################################################
# Configure instances
################################################################
variable "powervs_bootstrap_memory" {
  type        = string
  description = "Amount of memory, in  GiB, used by the bootstrap node."
  default     = "32"
}

variable "powervs_bootstrap_processors" {
  type        = string
  description = "Number of processors used by the bootstrap node."
  default     = "0.5"
}

variable "powervs_master_memory" {
  type        = string
  description = "Amount of memory, in  GiB, used by each master node."
  default     = "32"
}

variable "powervs_master_processors" {
  type        = string
  description = "Number of processors used by each master node."
  default     = "0.5"
}

variable "powervs_proc_type" {
  type        = string
  description = "The type of processor mode for all nodes (shared/dedicated)"
  default     = "shared"
}

variable "powervs_sys_type" {
  type        = string
  description = "The type of system (s922/e980)"
  default     = "s922"
}

variable "powervs_key_name" {
  type        = string
  description = "The name for the SSH key created in the Service Instance"
  default     = ""
}

variable "powervs_ssh_key" {
  type        = string
  description = "Public key for keypair used to access cluster. Required when creating 'ibm_pi_instance' resources."
  default     = ""
}

variable "powervs_image_bucket_name" {
  type        = string
  description = "Name of the COS bucket containing the image to be imported."
}

variable "powervs_image_bucket_file_name" {
  type        = string
  description = "File name of the image in the COS bucket."
}

variable "powervs_image_storage_type" {
  type        = string
  description = "Storage type used when storing image in Power VS."
  default     = "tier1"
}

variable "powervs_expose_bootstrap" {
  type        = bool
  description = "Setting this to false allows the bootstrap resources to be removed from the cluster load balancers."
  default     = true
}

################################################################
# Configure DNS
################################################################
## TODO: Pass the CIS CRN from the installer program, refer the IBM Cloud code to see the implementation.
variable "powervs_cis_crn" {
  type        = string
  description = "The CRN of CIS instance to use."
}

variable "powervs_dns_guid" {
  type        = string
  description = "The GUID of the IBM DNS Service instance to use when creating a private cluster."
}

################################################################
# Configure Tranit Gateways
################################################################

################################################################
# Output Variables
################################################################
variable "bootstrap_ip" { default = "" }
variable "control_plane_ips" { default = [] }
