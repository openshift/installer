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
  description = "Instance type for the bootstrap node. Example: `bx2-4x16`"
}

variable "ibmcloud_cis_crn" {
  type        = string
  description = "The CRN of CIS instance to use."
  default     = ""
}

variable "ibmcloud_dns_id" {
  type        = string
  description = "The ID of DNS Service instance to use."
  default     = ""
}

variable "ibmcloud_region" {
  type        = string
  description = "The target IBM Cloud region for the cluster."
}

variable "ibmcloud_master_instance_type" {
  type        = string
  description = "Instance type for the master node(s). Example: `bx2-4x16`"
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

variable "ibmcloud_terraform_private_visibility" {
  type        = bool
  description = "Specified whether the IBM Cloud terraform provider visibility mode should be private, for endpoint usage."
  default     = false
}

#######################################
# Top-level module variables (optional)
#######################################

variable "ibmcloud_endpoints_json_file" {
  type        = string
  description = "JSON file containing IBM Cloud service endpoints"
  default     = ""
}

variable "ibmcloud_preexisting_vpc" {
  type        = bool
  description = "Specifies whether an existing VPC should be used or a new one created for installation."
  default     = false
}

variable "ibmcloud_vpc_permitted" {
  type        = bool
  description = "Specifies whether an existing VPC is already a Permitted Network for DNS Instance, for Private clusters."
  default     = false
}

variable "ibmcloud_vpc" {
  type        = string
  description = "The name of an existing cluster VPC."
  default     = null
}

variable "ibmcloud_control_plane_boot_volume_key" {
  type        = string
  description = "IBM Cloud Key Protect key CRN to use to encrypt the control plane's volume(s)."
  default     = null
}

variable "ibmcloud_control_plane_subnets" {
  type        = list(string)
  description = "The names of the existing subnets for the control plane."
  default     = []
}

variable "ibmcloud_compute_subnets" {
  type        = list(string)
  description = "The names of the existing subnets for the compute plane."
  default     = []
}

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

variable "ibmcloud_network_resource_group_name" {
  type = string
  description = <<EOF
(optional) The name of the resource group for existing cluster network resources. If this is set, the existing network resources
(VPC, Subnets, etc.) must exist in the resource group to be used for cluster creation. Otherwise, new network resources are
created in the same resource group as the other cluster resources (see 'ibmcloud_resource_group_name').
EOF
  default = ""
}

variable "ibmcloud_resource_group_name" {
  type        = string
  description = <<EOF
(optional) The name of the resource group for the cluster. If this is set, the cluster is installed to that existing resource group
otherwise a new resource group will be created using cluster id.
EOF
  default = ""
}
