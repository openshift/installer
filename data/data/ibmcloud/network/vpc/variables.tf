#######################################
# VPC module variables
#######################################

variable "cluster_id" {
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