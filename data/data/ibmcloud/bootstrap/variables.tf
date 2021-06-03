#######################################
# Bootstrap module variables
#######################################

variable "cluster_id" {
  type = string
}

variable "cos_resource_instance_id" {
  type = string
}

variable "cos_bucket_region" {
  type = string
}

variable "ignition_file" {
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

variable "resource_group_id" {
  type = string
}

variable "security_group_id" {
  type = string
}

variable "subnet_id" {
  type = string
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

variable "zone" {
  type = string
}