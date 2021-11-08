variable "powervs_api_key" {
  type        = string
  description = "IBM Cloud API key associated with user's identity"
  default     = "<key>"
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
}

variable "cluster_id" {
    type    = string
    default = ""
}

variable "service_tags" {
  type        = list(string)
  description = "A list of tags for our resource instance."
  default     = []
}
