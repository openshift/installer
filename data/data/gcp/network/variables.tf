
variable "cluster_id" {
  type = string
}

variable "bootstrap_instances" {
  type        = list
  description = "The bootstrap instance."
}

variable "bootstrap_instance_groups" {
  type        = list
  description = "The bootstrap instance groups."
}

variable "bootstrap_lb" {
  type        = bool
  description = "If the bootstrap instance should be in the load balancers."
  default     = true
}

variable "master_instances" {
  type        = list
  description = "The master instances."
}

variable "master_instance_groups" {
  type        = list
  description = "The master instance groups."
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

variable "cluster_network" {
  type = string
}

variable "master_subnet" {
  type = string
}

variable "worker_subnet" {
  type = string
}

variable "preexisting_network" {
  type    = bool
  default = false
}

variable "public_endpoints" {
  type        = bool
  description = "If the bootstrap instance should have externally accessible resources."
}
