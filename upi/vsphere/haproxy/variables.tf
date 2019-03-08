variable "cluster" {
  type        = "string"
  description = "This is the vSphere server for the environment"
}

variable "datacenter" {
  type        = "string"
  description = "This is the vSphere server for the environment"
}

variable "datastore" {
  type        = "string"
  description = "This is the vSphere server for the environment"
}

variable "public_ipv4" {
  type        = "string"
  description = "This is the publicly accessibly endpoint."
}

variable "public_ipv4_gw" {
  type        = "string"
  description = "This is the public network gateway."
}

variable "public_netmask" {
  type        = "string"
  description = "This is the public network netmask."
}

variable "private_ipv4" {
  type        = "string"
  description = "This is the publicly accessibly endpoint."
}

variable "private_ipv4_gw" {
  type        = "string"
  description = "This is the public network gateway."
}

variable "private_netmask" {
  type        = "string"
  description = "This is the public network netmask."
}

variable "resource_pool" {
  type        = "string"
  description = "The resource pool for provisioning VMs."
}

variable "vm_network" {
  type        = "string"
  description = "This is the publicly accessible network for cluster ingress and access."
  default     = "VM Network"
}

variable "vm_private_network" {
  type        = "string"
  description = "This is the private network for cluster communication."
}

variable "vm_template" {
  type        = "string"
  description = "This is the RHEL template for environment use."
}
