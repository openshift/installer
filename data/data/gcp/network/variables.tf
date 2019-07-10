variable "bootstrap_instance" {
  type        = string
  description = "The bootstrap instance."
}

variable "bootstrap_instance_group" {
  type        = string
  description = "The instance group that hold the bootstrap instance in this region."
}

variable "cluster_id" {
  type = string
}

variable "worker_subnet_cidr" {
  type = string
}

variable "master_instances" {
  type        = list
  description = "The master instances."
}

variable "master_instance_groups" {
  type        = list
  description = "The instance groups that hold the master instances in this region."
}

variable "master_subnet_cidr" {
  type = string
}

variable "network_cidr" {
  type = string
}
