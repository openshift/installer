variable "etcd_ca_cert_pem" {
  type = "string"
}

variable "etcd_ca_key_alg" {
  type = "string"
}

variable "etcd_ca_key_pem" {
  type = "string"
}

variable "service_cidr" {
  type = "string"
}

variable "etcd_cert_dns_names" {
  type = "list"
}
