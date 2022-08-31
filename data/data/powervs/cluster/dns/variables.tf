variable "service_id" {
  type        = string
  description = "The ID of the IBM Cloud CIS instance, or IBM Cloud DNS instance, that will be used for the DNS records."
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

variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "vpc_crn" {
  type        = string
  description = "The CRN of the VPC."
}

variable "vpc_id" {
  type        = string
  description = "The ID of the VPC."
}

variable "vpc_subnet_id" {
  type        = string
  description = "The ID of the VPC subnet."
}

variable "vpc_zone" {
  type        = string
  description = "The IBM Cloud zone in which the VPC is created."
}

variable "dns_vm_image_name" {
  type        = string
  description = "The image name for the DNS VM."
  default     = "ibm-centos-7-9-minimal-amd64-6"
}

variable "ssh_key" {
  type        = string
  description = "Public key for keypair used to access cluster. Required when creating 'ibm_pi_instance' resources."
  default     = ""
}

variable "publish_strategy" {
  type        = string
  description = "The cluster publishing strategy, either Internal or External"
  default     = "External"
}

