################################################################
# Configure the IBM Cloud provider
################################################################
variable "ibmcloud_api_key" {
    type        = string
    description = "IBM Cloud API key associated with user's identity"
    default     = "<key>"

    validation{
        condition       = var.ibmcloud_api_key != "" && lower(var.ibmcloud_api_key) != "<key>"
        error_message   = "The ibmcloud_api_key is required and cannot be empty."
    }
}

variable "ibmcloud_region" {
    type        = string
    description = "The IBM Cloud region where you want to create the resources"
    default     = ""

    validation{
        condition       = var.ibmcloud_region != "" && lower(var.ibmcloud_region) != "<region>"
        error_message   = "The ibmcloud_region is required and cannot be empty."
    }
}

variable "ibmcloud_zone" {
    type        = string
    description = "The zone of an IBM Cloud region where you want to create Power System resources"
    default     = ""

    validation{
        condition       = var.ibmcloud_zone != "" && lower(var.ibmcloud_zone) != "<zone>"
        error_message   = "The ibmcloud_zone is required and cannot be empty."
    }
}

variable "cloud_instance_id" {
    type        = string
    description = "The cloud instance ID of your account"
    default     = ""
}

################################################################
# Configure instances
################################################################
variable "image_name" {
  type        = string
  description = "Name of the image used by all nodes in the cluster."
}

variable "network_name" {
  type        = string
  description = "Name of the network used by the all nodes in the cluster."
}

variable "bootstrap" {
    type    = object({memory = string, processors = string })
    default = {
        memory      = "32"
        processors  = "0.5"
    }
}

variable "proc_type" {
    type        = string
    description = "The type of processor mode for all nodes (shared/dedicated)"
    default     = "shared"
}

variable "sys_type" {
    type        = string
    description = "The type of system (s922/e980)"
    default     = "s922"
}

# Must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character
# Length cannot exceed 14 characters when combined with cluster_id_prefix
variable "cluster_id" {
    type    = string
    default = ""

    validation {
        condition     = can(regex("^$|^[a-z0-9]+[a-zA-Z0-9_\\-.]*[a-z0-9]+$", var.cluster_id))
        error_message = "The cluster_id value must be a lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character."
    }

    validation {
        condition     = length(var.cluster_id) <= 14
        error_message = "The cluster_id value shouldn't be greater than 14 characters."
    }
}
