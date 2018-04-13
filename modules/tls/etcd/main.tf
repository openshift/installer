# client keys
# These are used for "api server"-to-etcd and "etcd operator"-to-etcd client communication
resource "tls_private_key" "etcd_client" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "etcd_client" {
  key_algorithm   = "${tls_private_key.etcd_client.algorithm}"
  private_key_pem = "${tls_private_key.etcd_client.private_key_pem}"

  subject {
    common_name  = "etcd"
    organization = "etcd"
  }
}

resource "tls_locally_signed_cert" "etcd_client" {
  cert_request_pem = "${tls_cert_request.etcd_client.cert_request_pem}"

  ca_key_algorithm      = "${var.etcd_ca_key_alg}"
  ca_private_key_pem    = "${var.etcd_ca_key_pem}"
  ca_cert_pem           = "${var.etcd_ca_cert_pem}"
  validity_period_hours = "26280"

  allowed_uses = [
    "key_encipherment",
    "client_auth",
  ]
}
