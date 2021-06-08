variable "powervs_api_key" {
  type        = string
  description = "IBM Cloud API key associated with user's identity"
  default     = "<key>"

  validation{
    condition = var.powervs_api_key != "" && lower(var.powervs_api_key) != "<key>"
    error_message   = "The powervs_api_key is required and cannot be empty."
  }
}

variable "powervs_resource_group" {
  type        = string
  description = "The cloud instance resource group"
  default     = ""
}

variable "powervs_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = ""

  validation{
    condition       = var.powervs_region != "" && lower(var.powervs_region) != "<region>"
    error_message   = "The powervs_region is required and cannot be empty."
  }
}

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

variable "service_tags" {
  type        = list(string)
  description = "A list of tags for our resource instance."
  default     = []
}
