output "ca_cert_pem" {
  value = "${var.ca_cert_pem_path == "" ? var.ca_cert_pem : file(local._ca_cert_pem_path)}"
}

output "key_pem" {
  value = "${var.key_pem_path == "" ? join("", tls_private_key.ingress.*.private_key_pem) : file(local._key_pem_path)}"
}

output "cert_pem" {
  value = "${tls_locally_signed_cert.ingress.cert_pem}"
  value = "${var.cert_pem_path == "" ? join("", tls_locally_signed_cert.ingress.*.cert_pem) : file(local._cert_pem_path)}"
}
