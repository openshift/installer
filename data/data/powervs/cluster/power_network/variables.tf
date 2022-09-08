variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "cloud_instance_id" {
  type        = string
  description = "The Power VS Service Instance (aka Cloud Instance) ID."
}

variable "resource_group" {
  type        = string
  description = "The name of the Power VS resource group to which the user belongs."
}

variable "vpc_crn" {
  type        = string
  description = "The CRN of the IBM Cloud VPC."
}

variable "pvs_network_name" {
  type        = string
  description = "The name of a pre-created Power VS DHCP Network."
  default     = ""
}

variable "cloud_conn_name" {
  type        = string
  description = "The name of a pre-created Power VS Cloud connection."
  default     = ""
}

variable "machine_cidr" {
  type        = string
  description = "The machine network (IPv4 only)"
}
