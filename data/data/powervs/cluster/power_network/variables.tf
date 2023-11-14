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

variable "machine_cidr" {
  type        = string
  description = "The machine network (IPv4 only)"
}

variable "dns_server" {
  type        = string
  description = "The desired DNS server for the DHCP instance to server."
  default     = "1.1.1.1"
}

variable "enable_snat" {
  type        = bool
  description = "Boolean indicating if SNAT should be enabled."
  default     = true
}

variable "transit_gateway_enabled" {
  type        = bool
  description = "Boolean indicating if Transit Gateways should be used."
  default     = false
}
