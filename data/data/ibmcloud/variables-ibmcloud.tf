#######################################
# Top-level module variables (required)
#######################################

variable "ibmcloud_api_key" {
  type        = string
  # TODO: Supported on tf 0.14
  # sensitive   = true
  description = "The IAM API key for authenticating with IBM Cloud APIs."
}

variable "ibmcloud_cis_crn" {
  type        = string
  description = "The CRN of CIS instance to use."
}

variable "ibmcloud_region" {
  type        = string
  description = "The target IBM Cloud region for the cluster."
}

variable "ibmcloud_master_availability_zones" {
  type        = list(string)
  description = "The availability zones in which to create the masters. The length of this list must match master_count."
}

variable "ibmcloud_vsi_image" {
  type        = string
  description = "Name of VPC VSI image to use for all nodes."
}

#######################################
# Top-level module variables (optional)
#######################################

variable "ibmcloud_resource_group_name" {
  type    = string
  default = ""
}
variable "ibmcloud_vsi_profile" {
  type        = string
  description = "Name of VPC VSI profile to use for all nodes."
  default     = "bx2d-4x16"
}
