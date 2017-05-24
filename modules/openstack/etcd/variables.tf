// The content of the /etc/resolv.conf file.
variable resolv_conf_content {
  type = "string"
}

variable "base_domain" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "container_image" {
  type = "string"
}

variable core_public_keys {
  type = "list"
}

variable "tectonic_experimental" {
  default = false
}
