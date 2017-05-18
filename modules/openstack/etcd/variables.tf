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

variable "instance_count" {
  default = ""
}

variable "tls_enabled" {
  default = false
}

variable "tls_ca_crt_pem" {
  default = ""
}

variable "tls_client_key_pem" {
  default = ""
}

variable "tls_client_crt_pem" {
  default = ""
}

variable "tls_peer_key_pem" {
  default = ""
}

variable "tls_peer_crt_pem" {
  default = ""
}
