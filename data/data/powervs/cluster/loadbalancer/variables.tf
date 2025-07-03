variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "enable_snat" {
  type        = bool
  description = "Boolean indicating if SNAT should be enabled or disabled."
  default     = true
}

variable "master_count" {
  type        = string
  description = "The number of master nodes."
}

variable "vpc_id" {
  type        = string
  description = "The ID of the VPC."
}

variable "vpc_subnet_id" {
  type        = string
  description = "The ID of the VPC subnet."
}

variable "resource_group" {
  type        = string
  description = "The name of the Power VS resource group to which the user belongs."
}
