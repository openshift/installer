variable "vpc_id" {
  type = "string"
}

variable "tectonic_base_domain" {
  type = "string"
}

variable "tectonic_dns_name" {
  type = "string"
}

variable "console-elb" {
  type = "map"
}

variable "api-internal-elb" {
  type = "map"
}

variable "api-external-elb" {
  type = "map"
}
