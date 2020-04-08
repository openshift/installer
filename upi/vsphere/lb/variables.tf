variable "lb_ip_address" {
  type = string
}

variable "api_backend_addresses" {
  type = list(string)
}

variable "ingress_backend_addresses" {
  type = list(string)
}

variable "ssh_public_key_path" {
  type = string
}
