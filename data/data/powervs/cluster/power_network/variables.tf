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
