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
  description = "The contents of the ignition file."
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

variable "network_name" {
  type        = string
  description = "The name of the Power VS network."
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

variable "key_id" {
  type        = string
  description = "The SSH Key ID."
}

variable "cos_instance_location" {
  type        = string
  description = "Specifies whether the Cloud Object Store instance is global or in a specific region. Used for the ignition file."
}

variable "cos_bucket_location" {
  type        = string
  description = "The region in which to create the Cloud Object Store bucket. Used for the igntion file."
}

variable "cos_storage_class" {
  type        = string
  description = "The storage class for the Cloud Object Store instance."
}
