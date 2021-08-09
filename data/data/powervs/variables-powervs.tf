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
  ##default     = "eu-gb"
}

variable "powervs_vpc_zone" {
  type        = string
  description = "The zone of an IBM Cloud region where you want to create Power System resources"
  default     = "eu-gb-2"
}

variable "powervs_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = "lon"
}

variable "powervs_zone" {
  type        = string
  description = "The zone of an IBM Cloud region where you want to create Power System resources"
  default     = "lon04"
}

variable "powervs_resource_group" {
  type        = string
  description = "The cloud instance resource group"
  default     = ""
}

variable "powervs_cloud_instance_id" {
  type        = string
  description = "The cloud instance ID of your account"
  ## TODO: erase default and set via install-config
  default = "e449d86e-c3a0-4c07-959e-8557fdf55482"
}

################################################################
# Configure storage
################################################################
variable "powervs_cos_instance_location" {
  type        = string
  description = "The location of your COS instance"
  default     = "global"
}

variable "powervs_cos_bucket_location" {
  type        = string
  description = "The location to create your COS bucket"
  default     = "us-east"
}

variable "powervs_cos_storage_class" {
  type        = string
  description = "The plan used for your COS instance"
  default     = "smart"
}

################################################################
# Configure instances
################################################################
variable "powervs_image_name" {
  type        = string
  description = "Name of the image used by all nodes in the cluster."
}

variable "powervs_network_name" {
  type        = string
  description = "Name of the network used by the all nodes in the cluster."
  default     = "pvs-ipi-net"
}

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

variable "powervs_ssh_key" {
  type        = string
  description = "Public key for keypair used to access cluster. Required when creating 'ibm_pi_instance' resources."
}

## TODO: Set this in install-config instead
variable "powervs_vpc_name" {
  type        = string
  description = "Name of the IBM Cloud Virtual Private Cloud (VPC) to setup the load balancer."
  default     = "powervs-ipi"
}

variable "powervs_vpc_subnet_name" {
  type        = string
  description = "Name of the VPC subnet having DirectLink access to the PowerVS private network"
  default     = "subnet2"
}

## TODO: Pass the CIS CRN from the installer program, refer the IBM Cloud code to see the implementation.
variable "powervs_cis_crn" {
  type        = string
  description = "The CRN of CIS instance to use."
  default     = "crn:v1:bluemix:public:internet-svcs:global:a/65b64c1f1c29460e8c2e4bbfbd893c2c:453c4cff-2ee0-4309-95f1-2e9384d9bb96::"
}
