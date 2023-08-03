variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "resource_group" {
  type        = string
  description = "The name of the Power VS resource group to which the user belongs."
}

variable "transit_gateway_enabled" {
  type        = bool
  description = "Boolean indicating if Transit Gateways should be used."
  default     = false
}

variable "service_instance_crn" {
  type        = string
  description = "The CRN of the service instance."
}

variable "vpc_crn" {
  type        = string
  description = "The CRN of the IBM Cloud VPC."
}

variable "vpc_region" {
  type        = string
  description = "The IBM Cloud region in which the VPC is created."
}
