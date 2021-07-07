################################################################
# Configure the IBM Cloud provider
################################################################
variable "ibmcloud_api_key" {
  type        = string
  description = "IBM Cloud API key associated with user's identity"
  default     = "<key>"

  validation {
    condition     = var.ibmcloud_api_key != "" && lower(var.ibmcloud_api_key) != "<key>"
    error_message = "The ibmcloud_api_key is required and cannot be empty."
  }
}

variable "ibmcloud_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = ""

  validation {
    condition     = var.ibmcloud_region != "" && lower(var.ibmcloud_region) != "<region>"
    error_message = "The ibmcloud_region is required and cannot be empty."
  }
}

variable "ibmcloud_zone" {
  type        = string
  description = "The zone of an IBM Cloud region where you want to create Power System resources"
  default     = ""

  validation {
    condition     = var.ibmcloud_zone != "" && lower(var.ibmcloud_zone) != "<zone>"
    error_message = "The ibmcloud_zone is required and cannot be empty."
  }
}

variable "powervs_resource_group" {
  type        = string
  description = "The cloud instance resource group"
  default     = ""
}

variable "powervs_cloud_instance_id" {
  type        = string
  description = "The cloud instance ID of your account"
  default     = ""
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

variable "powervs_base_domain" {
  type        = string
  description = "The base domain name of the cluster"
  default     = ""
}

variable "powervs_cluster_domain" {
  type        = string
  description = "The name of the cluster that all DNS records must belong to."
  default     = ""
}

variable "powervs_ssh_key" {
  type        = string
  description = "Public key for keypair used to access cluster."
  default     = ""
}

variable "powervs_vpc_name" {
  type        = string
  description = "Name of the IBM Cloud Virtual Private Cloud (VPC) to setup the load balancer."
  default     = ""
}

variable "powervs_vpc_subnet_name" {
  type        = string
  description = "Name of the VPC subnet having DirectLink access to the PowerVS private network"
  default     = ""
}
