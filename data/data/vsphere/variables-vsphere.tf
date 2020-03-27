//////
// vSphere variables
//////

variable "vsphere_url" {
  type        = string
  description = "This is the vSphere server for the environment."
}

variable "vsphere_username" {
  type        = string
  description = "vSphere server user for the environment."
}

variable "vsphere_password" {
  type        = string
  description = "vSphere server password"
}

variable "vsphere_cluster" {
  type        = string
  description = "This is the name of the vSphere cluster."
}

variable "vsphere_datacenter" {
  type        = string
  description = "This is the name of the vSphere data center."
}

variable "vsphere_datastore" {
  type        = string
  description = "This is the name of the vSphere data store."
}

variable "vsphere_ova_filepath" {
  type        = string
  description = "This is the filepath to the ova file that will be imported into vSphere."
}

variable "vsphere_template" {
  type        = string
  description = "This is the name of the VM template to clone."
}

variable "vsphere_network" {
  type        = string
  description = "This is the name of the publicly accessible network for cluster ingress and access."
}

variable "vsphere_folder" {
  type = string
}

///////////
// Control Plane machine variables
///////////

variable "vsphere_control_plane_memory_mib" {
  type = number
}

variable "vsphere_control_plane_disk_gib" {
  type = number
}

variable "vsphere_control_plane_num_cpus" {
  type = number
}

variable "vsphere_control_plane_cores_per_socket" {
  type = number
}

