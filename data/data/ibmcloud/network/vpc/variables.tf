#######################################
# VPC module variables
#######################################

variable "cluster_id" {
  type = string
}

variable "network_resource_group_id" {
  type = string
}

variable "public_endpoints" {
  type = bool
}

variable "resource_group_id" {
  type = string
}

variable "tags" {
  type = list(string)
}

variable "zones_master" {
  type = list(string)
}

variable "zones_worker" {
  type = list(string)
}

variable "preexisting_vpc" {
  type    = bool
  default = false
}

variable "cluster_vpc" {
  type = string
}

variable "control_plane_subnets" {
  type = list(string)
}

variable "compute_subnets" {
  type = list(string)
}
