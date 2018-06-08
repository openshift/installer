# Admin (generated/tls/{admin.key,admin.crt})
# Used to create kubeconfig (generated/auth/kubeconfig) with admin level privileges.
resource "tls_private_key" "admin" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "admin" {
  key_algorithm   = "${tls_private_key.admin.algorithm}"
  private_key_pem = "${tls_private_key.admin.private_key_pem}"

  subject {
    common_name  = "system:admin"
    organization = "system:masters"
  }
}

resource "tls_locally_signed_cert" "admin" {
  cert_request_pem = "${tls_cert_request.admin.cert_request_pem}"

  ca_key_algorithm      = "${var.kube_ca_key_alg}"
  ca_private_key_pem    = "${var.kube_ca_key_pem}"
  ca_cert_pem           = "${var.kube_ca_cert_pem}"
  validity_period_hours = "26280"

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}
