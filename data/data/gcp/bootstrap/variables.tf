variable "network" {
  type = string
}

variable "master_subnet" {
  type = string
}

variable "compute_image" {
  type = string
}

variable "cluster_ip" {
  type = string
}

variable "cluster_public_ip" {
  type        = string
  default     = null
  description = "IP of the API load balancer; it is null with the internal publishing strategy."
}
