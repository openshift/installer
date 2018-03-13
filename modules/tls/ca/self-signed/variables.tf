variable "root_ca_cert_pem" {
  description = "PEM-encoded CA certificate (generated if blank)"
  type        = "string"
}

variable "root_ca_key_alg" {
  description = "Algorithm used to generate root_ca_key (required if root_ca_cert is specified)"
  type        = "string"
}

variable "root_ca_key_pem" {
  description = "PEM-encoded CA key (required if root_ca_cert is specified)"
  type        = "string"
}
