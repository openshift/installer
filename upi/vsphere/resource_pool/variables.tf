variable "vsphere_cluster" {
  type        = "string"
  description = "This is the vSphere server for the environment."
}

variable "vsphere_datacenter" {
  type        = "string"
  description = "This is the vSphere server for the environment."
}

variable "vsphere_resource_pool" {
  type        = "string"
  description = "The resource pool for provisioning VMs."
}
