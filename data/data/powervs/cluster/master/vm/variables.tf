variable "instance_count" {
  type        = string
  description = "The number of master nodes to be created."
}

variable "memory" {
  type        = string
  description = "The amount of memory to assign to each node in GB."
}

variable "processors" {
  type        = string
  description = "The processor count for each node."
}

variable "ignition" {
  type        = string
  description = "The ignition file."
}

variable "cloud_instance_id" {
  type        = string
  description = "The Power VS Service Instance (aka Cloud Instance) ID."
}

variable "resource_group" {
  type        = string
  description = "The name of the Power VS resource group to which the user belongs."
}

variable "image_id" {
  type        = string
  description = "The ID of the Power VS boot image for the nodes."
}

variable "proc_type" {
  type        = string
  description = "The type of processor to be assigned (e.g. capped, dedicated, shared) to the nodes."
}

variable "sys_type" {
  type        = string
  description = "The type of system on which to provision the nodes (e.g s922)."
}
variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "ssh_key_name" {
  type        = string
  description = "The SSH Key Name."
}

variable "dhcp_id" {
  type        = string
  description = "The ID of the Power VS DHCP Service."
}

variable "dhcp_network_id" {
  type        = string
  description = "The ID of the Power VS DHCP network."
}
