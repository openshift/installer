// usually the caller sets this to:
// var.tectonic_experimental || var.tectonic_etcd_tls_enabled
variable "self_signed" {
  description = <<EOF
If set to true, self-signed certificates are generated.
If set to false, only the passed CA and client certs are being used.
EOF
}

variable "etcd_ca_cert_path" {
  type        = "string"
  description = "external CA certificate"
}

variable "etcd_client_cert_path" {
  type = "string"
}

variable "etcd_client_key_path" {
  type = "string"
}

variable "service_cidr" {
  type = "string"
}

variable "etcd_cert_dns_names" {
  type = "list"
}
