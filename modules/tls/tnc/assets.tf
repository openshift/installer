resource "local_file" "tnc_cert" {
  content  = "${tls_locally_signed_cert.tnc.cert_pem}"
  filename = "./generated/tls/tnc.crt"
}

resource "local_file" "tnc_key" {
  content  = "${tls_private_key.tnc.private_key_pem}"
  filename = "./generated/tls/tnc.key"
}
