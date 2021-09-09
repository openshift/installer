variable "control_plane_ips" {
  type    = list(string)
  default = null
}

variable "lb_target_group_arns" {
  type = list(string)
}

variable "lb_target_group_arns_length" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "public_subnet_ids" {
  type = list(string)
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "master_sg_id" {
  type = string
}

variable "ami_id" {
  type = string
}

variable "bootstrap_ip" {
  type    = string
  default = null
}
