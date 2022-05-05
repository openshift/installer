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
  default     = ""
}

variable "ovirt_ca_bundle" {
  type        = string
  description = "The CA certificate for the oVirt engine API in PEM format"
}

variable "ovirt_insecure" {
  type        = bool
  description = "Disable oVirt engine certificate verification"
}

variable "ovirt_cluster_id" {
  type        = string
  description = "The ID of Cluster"
  validation {
    condition     = var.ovirt_cluster_id != ""
    error_message = "The ovirt_storage_domain_id must not be empty."
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

variable "ovirt_base_image_name" {
  type        = string
  default     = ""
  description = "Name of an existing base image to use for the nodes."
}

variable "ovirt_base_image_local_file_path" {
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

variable "ovirt_master_auto_pinning_policy" {
  type        = string
  default     = "none"
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

