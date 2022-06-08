variable "bootstrap_ip" {
  type        = string
  description = "The IP address of the bootstrap node."
}

variable "lb_ext_id" {
  type        = string
  description = "The ID of the external load balancer."
}

variable "lb_int_id" {
  type        = string
  description = "The ID of the internal load balancer."
}

variable "machine_cfg_pool_id" {
  type        = string
  description = "The ID of the machine config load balancer pool."
}

variable "api_pool_int_id" {
  type        = string
  description = "The ID of the internal API load balancer pool."
}

variable "api_pool_ext_id" {
  type        = string
  description = "The ID of the external API load balancer pool."
}
