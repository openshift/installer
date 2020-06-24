variable "cluster_id" {
  description = "The ID of Openshift cluster"
}

variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of Cluster"
}

variable "ovirt_storage_domain_id" {
  type        = string
  description = "The ID of Storage Domain"
}

variable "ignition_bootstrap" {
  type        = string
  description = "bootstrap ignition config"
}

variable "openstack_base_image_name" {
  type        = string
  description = "Name of the base image to use for the nodes."
}

variable "openstack_base_image_local_file_path" {
  type        = string
  default     = ""
  description = "Local file path of the base image file to use for the nodes."
}

variable "ovirt_network_name" {
  type        = string
  description = "The name of Logical Network for the selected Cluster."
}

variable "ovirt_vnic_profile_id" {
  type        = string
  description = "The ID of the vNIC profile of Logical Network."
}
