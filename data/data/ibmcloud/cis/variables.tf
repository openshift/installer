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

variable "bootstrap_name" {
  type = string
}

variable "bootstrap_ipv4_address" {
  type = string
}

variable "master_count" {
  type = string
}

variable "master_name_list" {
  type = list(string)
}

variable "master_ipv4_address_list" {
  type = list(string)
}

variable "lb_kubernetes_api_public_hostname" {
  type = string
}

variable "lb_kubernetes_api_private_hostname" {
  type = string
}
