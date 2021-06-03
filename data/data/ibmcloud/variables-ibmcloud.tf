#######################################
# Top-level module variables (required)
#######################################

variable "ibmcloud_api_key" {
  type        = string
  # TODO: Supported on tf 0.14
  # sensitive   = true
  description = "The IAM API key for authenticating with IBM Cloud APIs."
}

variable "ibmcloud_bootstrap_instance_type" {
  type        = string
  description = "Instance type for the bootstrap node. Example: `bx2d-4x16`"
}

variable "ibmcloud_cis_crn" {
  type        = string
  description = "The CRN of CIS instance to use."
}

variable "ibmcloud_region" {
  type        = string
  description = "The target IBM Cloud region for the cluster."
}

variable "ibmcloud_master_instance_type" {
  type        = string
  description = "Instance type for the master node(s). Example: `bx2d-4x16`"
}

variable "ibmcloud_master_availability_zones" {
  type        = list(string)
  description = "The availability zones in which to create the masters. The length of this list must match master_count."
}

variable "ibmcloud_image_filepath" {
  type        = string
  description = "The file path to the RHCOS image"
}

#######################################
# Top-level module variables (optional)
#######################################

variable "ibmcloud_resource_group_name" {
  type    = string
  default = ""
}
