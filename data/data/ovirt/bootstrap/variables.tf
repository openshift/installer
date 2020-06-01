variable "cluster_id" {
  description = "The ID of Openshift cluster"
}

variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of Cluster"
}

variable "ovirt_template_id" {
  type        = string
  description = "The ID of VM template"
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
