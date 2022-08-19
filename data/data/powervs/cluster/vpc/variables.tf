variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "resource_group" {
  type        = string
  description = "The name of the Power VS resource group to which the user belongs."
}

variable "vpc_zone" {
  type        = string
  description = "The IBM Cloud zone in which the VPC is created."
}

variable "wait_for_vpc" {
  type        = string
  description = "The seconds wait for VPC creation, default is 60s."
}

variable "vpc_subnet_name" {
  type        = string
  description = "The name of a pre-created VPC subnet."
  default     = ""
}

variable "vpc_name" {
  type        = string
  description = "The name of a pre-created VPC."
  default     = ""
}

