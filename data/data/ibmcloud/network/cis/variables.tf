############################################
# CIS module variables
############################################

variable "cis_id" {
  type = string
}

variable "base_domain" {
  type = string
}

variable "cluster_domain" {
  type = string
}

variable "lb_kubernetes_api_public_hostname" {
  type = string
}

variable "lb_kubernetes_api_private_hostname" {
  type = string
}
