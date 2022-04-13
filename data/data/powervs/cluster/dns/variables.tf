variable "cis_id" {
  type        = string
  description = "The ID of the IBM Cloud CIS instance that will be used for the DNS records."
}

variable "base_domain" {
  type        = string
  description = "The base domain for all DNS records."
}

variable "cluster_domain" {
  type        = string
  description = "The domain name for the created cluster."
}

variable "load_balancer_hostname" {
  type        = string
  description = "The hostname for the external load balancer."
}

variable "load_balancer_int_hostname" {
  type        = string
  description = "The hostname for the internal load balancer."
}




