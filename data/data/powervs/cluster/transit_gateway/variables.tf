variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "resource_group" {
  type        = string
  description = "The name of the Power VS resource group to which the user belongs."
}

variable "attached_transit_gateway" {
  type        = string
  description = "The ID of already attached Transit Gateways, if any."
}

variable "tg_connection_vpc_id" {
  type        = string
  description = "ID of a VPC connection to the transit gateway specified in attached_transit_gateway, if any."
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
