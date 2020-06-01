variable "api_server" {
  type = map
  default = {
    "lb_port" = "6443"
    "lb_protocol" = "TCP"
  }
}

variable "machine_config_server" {
  type = map
  default = { 
    "lb_port" = "22623"
    "lb_protocol" = "TCP"
  }
}

variable "ingress_http" {
  type = map
  default = {
    "lb_port" = "80"
    "lb_protocol" = "TCP"
  }
}

variable "ingress_https" {
  type = map
  default = {
    "lb_port" = "443"
    "lb_protocol" = "TCP"
  }
}

variable "lb_ip_address" {
  type = string
}

variable "api_backend_addresses" {
  type = list(string)
}

variable "ingress_backend_addresses" {
  type = list(string)
}

variable "citrix_adc_username" {
  type = string
}

variable "citrix_adc_password" {
  type = string
}

variable "citrix_adc_ip" {
  type = string
}
