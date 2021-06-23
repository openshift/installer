
variable "cluster_id" {
  type = string
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
