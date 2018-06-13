# These are used for Ignition-to-TNC communication
resource "tls_private_key" "tnc" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "tnc" {
  key_algorithm   = "${tls_private_key.tnc.algorithm}"
  private_key_pem = "${tls_private_key.tnc.private_key_pem}"

  subject {
    common_name = "${var.domain}"
  }

  dns_names = [
    "${var.domain}",
  ]
}

resource "tls_locally_signed_cert" "tnc" {
  cert_request_pem = "${tls_cert_request.tnc.cert_request_pem}"

  ca_key_algorithm      = "${var.ca_key_alg}"
  ca_private_key_pem    = "${var.ca_key_pem}"
  ca_cert_pem           = "${var.ca_cert_pem}"
  validity_period_hours = "26280"

  allowed_uses = [
    "server_auth",
  ]
}
