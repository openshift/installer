variable "master_ips" {
  type        = list(any)
  description = "The IP addresses of the master nodes."
}

variable "instance_count" {
  type        = string
  description = "The number of master nodes to be created."
}

variable "lb_ext_id" {
  type        = string
  description = "The ID of the external load balancer in the IBM Cloud VPC"
}

variable "lb_int_id" {
  type        = string
  description = "The ID of the private load balancer in the IBM Cloud VPC"
}

variable "machine_cfg_pool_id" {
  type        = string
  description = "The ID of the load balancer pool for the machine-config server."
}

variable "api_pool_int_id" {
  type        = string
  description = "The ID of the load balancer pool for the API server."
}

variable "api_pool_ext_id" {
  type        = string
  description = "The ID of the public load balancer pool for the API server."
}

