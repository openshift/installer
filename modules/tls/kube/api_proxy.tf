# Kubernetes API Server Proxy (resources/generated/tls/{apiserver-proxy.key,apiserver-proxy.crt})
resource "tls_private_key" "apiserver_proxy" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "apiserver_proxy" {
  key_algorithm   = "${tls_private_key.apiserver_proxy.algorithm}"
  private_key_pem = "${tls_private_key.apiserver_proxy.private_key_pem}"

  subject {
    common_name  = "kube-apiserver-proxy"
    organization = "kube-master"
  }
}

resource "tls_locally_signed_cert" "apiserver_proxy" {
  cert_request_pem = "${tls_cert_request.apiserver_proxy.cert_request_pem}"

  ca_key_algorithm      = "${var.aggregator_ca_key_alg}"
  ca_private_key_pem    = "${var.aggregator_ca_key_pem}"
  ca_cert_pem           = "${var.aggregator_ca_cert_pem}"
  validity_period_hours = "26280"

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}
