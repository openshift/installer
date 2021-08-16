#######################################
# Master module variables
#######################################

variable "control_plane_security_group_id_list" {
  type = list(string)
}

variable "control_plane_subnet_id_list" {
  type = list(string)
}

variable "control_plane_subnet_zone_list" {
  type = list(string)
}

variable "cos_resource_instance_id" {
  type    = string
  default = ""
}

variable "lb_kubernetes_api_public_id" {
  type = string
}

variable "lb_kubernetes_api_private_id" {
  type = string
}

variable "lb_pool_kubernetes_api_public_id" {
  type = string
}

variable "lb_pool_kubernetes_api_private_id" {
  type = string
}

variable "lb_pool_machine_config_id" {
  type = string
}

variable "resource_group_id" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "vsi_image_id" {
  type = string
}
