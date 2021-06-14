variable "bootstrap_instances" {
  type        = list
  description = "The bootstrap instance."
}

variable "bootstrap_instance_groups" {
  type        = list
  description = "The bootstrap instance groups."
}

variable "master_instances" {
  type        = list
  description = "The master instances."
}

variable "master_instance_groups" {
  type        = list
  description = "The master instance groups."
}

variable "network" {
  type = string
}

variable "master_subnet" {
  type = string
}

variable "cluster_ip" {
  type = string
}

variable "cluster_public_ip" {
  type = string
}

variable "api_health_checks" {
  type = list
}

variable "api_internal_health_checks" {
  type = list
}