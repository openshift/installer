resource "local_file" "ca_cert" {
  content  = "${var.ca_cert_pem_path == "" ? var.ca_cert_pem : file(local._ca_cert_pem_path)}"
  filename = "./generated/tls/ingress-ca.crt"
}

resource "local_file" "cert" {
  content  = "${var.cert_pem_path == "" ? join("", tls_locally_signed_cert.ingress.*.cert_pem) : file(local._cert_pem_path)}"
  filename = "./generated/tls/ingress.crt"
}

resource "local_file" "key" {
  content  = "${var.key_pem_path == "" ? join("", tls_private_key.ingress.*.private_key_pem) : file(local._key_pem_path)}"
  filename = "./generated/tls/ingress.key"
}
