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
  description = "The ID of Cluster"
}

variable "ovirt_template_id" {
  type        = string
  description = "The ID of VM template"
}

variable "ignition_master" {
  type        = string
  description = "master ignition config"
}

variable "ovirt_master_memory" {
  type        = string
  description = "master VM memory in MiB"
  default     = 16348 * 1024 * 1024
}

variable "ovirt_master_cores" {
  type        = string
  description = "master VM number of cores"
  default     = 1
}

variable "ovirt_master_sockets" {
  type        = string
  description = "master VM number of sockets"
  default     = 1
}

variable "ovirt_master_threads" {
  type        = string
  description = "master VM number of threads"
  default     = 1
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

variable "ovirt_master_affinity_groups" {
  type        = list(string)
  description = "master VMs affinity groups names"
  default     = []
}

//TODO: REMOVE once we port to TF 0.13 and can use depends_on modules
variable "ovirt_affinity_group_count" {
  type        = string
  description = "create a dependency between affinity_group module to masters module"
}

variable "ovirt_master_auto_pinning_policy" {
  type        = string
  description = "master VM auto pinning policy"
}

variable "ovirt_master_hugepages" {
  type        = string
  description = "master VM hugepages size in KiB"
}

variable "ovirt_master_sparse" {
  type        = bool
  description = "make master VM disks sparse."
  default     = null
}

variable "ovirt_master_clone" {
  type        = bool
  description = "clone master VM disk from template instead of linking. Defaults to false for desktop ovirt_master_vm_type, true otherwise."
  default     = null
}

variable "ovirt_master_format" {
  type        = string
  description = "master VM disk format, can be empty, 'raw', or 'cow'"
  validation {
    condition     = var.ovirt_master_format == "" || var.ovirt_master_format == "cow" || var.ovirt_master_format == "raw"
    error_message = "The ovirt_master_format must be one of 'raw' or 'cow'."
  }
}

variable "ovirt_storage_domain_id" {
  type        = string
  description = "The ID of Storage Domain for the template"
  validation {
    condition     = var.ovirt_storage_domain_id != ""
    error_message = "The ovirt_storage_domain_id must not be empty."
  }
}