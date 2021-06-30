#######################################
# Master module variables
#######################################

variable "cluster_id" {
  type = string
}

variable "instance_count" {
  type = string
}

variable "ignition" {
  type = string
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

variable "public_endpoints" {
  type = bool
}

variable "resource_group_id" {
  type = string
}

variable "security_group_id_list" {
  type = list(string)
}

variable "subnet_id_list" {
  type = list(string)
}

variable "tags" {
  type = list(string)
}

variable "vpc_id" {
  type = string
}

variable "vsi_image_id" {
  type = string
}

variable "vsi_profile" {
  type = string
}

variable "zone_list" {
  type = list(string)
}