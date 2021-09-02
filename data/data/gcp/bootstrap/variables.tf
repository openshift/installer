variable "cluster_ip" {
  type    = string
  default = null
}

variable "cluster_public_ip" {
  type    = string
  default = null
}

variable "network" {
  type = string
}

variable "master_subnet" {
  type = string
}

variable "api_health_checks" {
  type = list
}

variable "api_internal_health_checks" {
  type = list
}

variable "master_instances" {
  type = list
}

variable "master_instance_groups" {
  type = list
}

variable "compute_image" {
  type = string
}

variable "control_plane_ips" {
  type = list
}
