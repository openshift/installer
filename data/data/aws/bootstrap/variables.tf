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

variable "edge_public_subnet_ids" {
  type = list(string)
}

variable "edge_private_subnet_ids" {
  type = list(string)
}

variable "master_sg_id" {
  type = string
}

variable "ami_id" {
  type = string
}

variable "aws_external_api_lb_dns_name" {
  type = string
}

variable "aws_internal_api_lb_dns_name" {
  type = string
}
