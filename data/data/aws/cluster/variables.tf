variable "aws_lb_api_external_dns_name" {
  type = string
}

variable "aws_lb_api_external_zone_id" {
  type = string
}

variable "aws_lb_api_internal_dns_name" {
  type = string
}

variable "aws_lb_api_internal_zone_id" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "az_to_private_subnet_id" {
  type = map(string)
}

variable "master_sg_id" {
  type = string
}

variable "master_sg_ids" {
  type = list(string)
}

variable "aws_lb_target_group_arns" {
  type = list(string)
}

variable "aws_lb_target_group_arns_length" {
  type = string
}

variable "ami_id" {
  type = string
}
