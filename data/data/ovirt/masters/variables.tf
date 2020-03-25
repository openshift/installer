variable "cluster_id" {
  description = "The ID of Openshift cluster"
}

variable "cluster_domain" {
  description = "The domain name of Openshift cluster"
}

variable "master_count" {
  type        = string
  description = "Number of masters"
  default     = 3
}

variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of oVirt's cluster"
}

variable "ovirt_template_id" {
  type        = string
  description = "The ID of oVirt's VM template"
}

variable "ignition_master" {
  type        = string
  description = "master ignition config"
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

variable "ovirt_master_os_disk_size_gb" {
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
