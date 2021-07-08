variable "cluster_id" {
  description = "The ID of Openshift cluster"
}

variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of Cluster"
}

variable "openstack_base_image_name" {
  type        = string
  description = "Name of the existing base image to use for the nodes."
}

variable "tmp_import_vm_id" {
  type        = string
  description = "ID of the temporary VM template created"
}
