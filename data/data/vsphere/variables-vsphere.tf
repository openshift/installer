variable "vsphere_disk_type" {
  type = string
}

variable "vsphere_vcenters" {
  type    = map(any)
  default = {}
}

variable "vsphere_networks" {
  type    = map(any)
  default = {}
}

variable "vsphere_folders" {
  type    = map(any)
  default = {}
}

variable "vsphere_control_planes" {
  type    = list(any)
  default = []
}

variable "vsphere_failure_domains" {
  type    = list(any)
  default = []
}

variable "vsphere_ova_filepath" {
  type        = string
  description = "This is the filepath to the ova file that will be imported into vSphere."
}
