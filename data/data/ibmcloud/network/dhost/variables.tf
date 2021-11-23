#######################################
# Dedicated Host module variables
#######################################

variable "cluster_id" {
  type = string
}

variable "dedicated_hosts_master" {
  type    = list(map(string))
  default = []
}

variable "dedicated_hosts_worker" {
  type    = list(map(string))
  default = []
}

variable "resource_group_id" {
  type = string
}

variable "zones_master" {
  type = list(string)
}

variable "zones_worker" {
  type = list(string)
}