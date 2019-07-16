variable "bootstrap_instances" {
  type        = list
  description = "The bootstrap instance."
}

variable "bootstrap_instance_groups" {
  type        = list
  description = "The instance group that hold the bootstrap instance in this region."
}

variable "bootstrap_lb" {
  type        = bool
  description = "If the bootstrap instance should be in the load balancers."
  default     = true
}

variable "cluster_id" {
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

variable "worker_subnet_cidr" {
  type = string
}
