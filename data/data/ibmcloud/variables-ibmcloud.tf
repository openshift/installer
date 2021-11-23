#######################################
# Top-level module variables (required)
#######################################

variable "ibmcloud_api_key" {
  type = string
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

variable "ibmcloud_worker_availability_zones" {
  type        = list(string)
  description = "The availability zones to provision for workers. Worker instances are created by the machine-API operator, but this variable controls their supporting infrastructure (subnets, routing, dedicated hosts, etc.)."
}

variable "ibmcloud_image_filepath" {
  type        = string
  description = "The file path to the RHCOS image"
}

#######################################
# Top-level module variables (optional)
#######################################

variable "ibmcloud_master_dedicated_hosts" {
  type        = list(map(string))
  description = "(optional) The list of dedicated hosts in which to create the control plane nodes."
  default     = []
}

variable "ibmcloud_worker_dedicated_hosts" {
  type        = list(map(string))
  description = "(optional) The list of dedicated hosts in which to create the compute nodes."
  default     = []
}

variable "ibmcloud_extra_tags" {
  type        = list(string)
  description = <<EOF
(optional) Extra IBM Cloud tags to be applied to created resources.
Example: `[ "key:value", "foo:bar" ]`
EOF
  default = []
}

variable "ibmcloud_publish_strategy" {
  type = string
  description = "The cluster publishing strategy, either Internal or External"
  default = "External"
  # TODO: Supported on tf 0.13
  # validation {
  #   condition     = "External" || "Internal"
  #   error_message = "The ibmcloud_publish_strategy value must be \"External\" or \"Internal\"."
  # }
}

variable "ibmcloud_resource_group_name" {
  type = string
  description = <<EOF
(optional) The name of the resource group for the cluster. If this is set, the cluster is installed to that existing resource group
otherwise a new resource group will be created using cluster id.
EOF
  default = ""
}
