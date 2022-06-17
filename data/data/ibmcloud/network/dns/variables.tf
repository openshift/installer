############################################
# DNS module variables
############################################

variable "dns_id" {
  type = string
}

variable "vpc_crn" {
  type = string
}

variable "vpc_permitted" {
  type = bool
}

variable "base_domain" {
  type = string
}

variable "cluster_domain" {
  type = string
}

variable "is_external" {
  type = bool
}

variable "lb_kubernetes_api_private_hostname" {
  type = string
}
