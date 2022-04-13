variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "vpc_name" {
  type        = string
  description = "The name of the VPC if precreated."
}

variable "vpc_subnet_id" {
  type        = string
  description = "The name of the VPC subnet if precreated."
}

variable "bootstrap_ip" {
  type        = string
  description = "The IP address of the bootstrap node."
}

variable "master_ips" {
  type        = list(string)
  description = "The IP addresses of the master nodes."
}

variable "resource_group" {
  type        = string
  description = "The name of the Power VS resource group to which the user belongs."
}
