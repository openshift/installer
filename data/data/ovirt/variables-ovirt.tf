variable "bootstrap_dns" {
  type        = string
  default     = true
  description = "Whether to include DNS entries for the bootstrap node or not."
}

variable "ovirt_url" {
  type        = string
  description = "The Engine URL"
}

variable "ovirt_username" {
  type        = string
  description = "The name of user to access Engine API"
}

variable "ovirt_password" {
  type        = string
  description = "The plain password of user to access Engine API"
}

variable "ovirt_cafile" {
  type        = string
  description = "Path to a file containing the CA certificate for the oVirt engine API in PEM format"
}

variable "ovirt_ca_bundle" {
  type        = string
  description = "The CA certificate for the oVirt engine API in PEM format"
}

variable "ovirt_insecure" {
  type        = bool
  description = "Disable oVirt engine certificate verification"
}

variable "ovirt_tmp_template_vm_id" {
  type        = string
  default     = ""
  description = "The ID of tmp VM which was created for creating the templated"
}

variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of Cluster"
}

variable "ovirt_storage_domain_id" {
  type        = string
  description = "The ID of Storage Domain for the template"
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
  default     = "ovirtmgmt"
  description = "The name of Logical Network for the selected Engine cluster."
}

variable "ovirt_vnic_profile_id" {
  type        = string
  description = "The ID of the vNIC profile of Logical Network."
}

variable "ovirt_affinity_groups" {
  type        = list(object({ name = string, priority = number, description = string, enforcing = string }))
  description = "Affinity groups that will be created"
  default     = []
}

variable "ovirt_master_memory" {
  type        = string
  description = "master VM memory in MiB"
}

variable "ovirt_master_cores" {
  type        = string
  description = "master VM number of cores"
}

variable "ovirt_master_sockets" {
  type        = string
  description = "master VM number of sockets"
}

variable "ovirt_master_os_disk_gb" {
  type        = string
  description = "master VM disk size in GiB"
}

variable "ovirt_master_vm_type" {
  type        = string
  description = "master VM type"
}

variable "ovirt_master_instance_type_id" {
  type        = string
  description = "master VM instance type ID"
}

variable "ovirt_master_affinity_groups" {
  type        = list(string)
  description = "master VMs affinity groups names"
}
