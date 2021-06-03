#######################################
# VPC module variables
#######################################

variable "cluster_id" {
  type = string
}

variable "region" {
  type = string
}

variable "resource_group_id" {
  type = string
}

variable "tags" {
  type = list(string)
}

variable "zone_list" {
  type = list(string)
}