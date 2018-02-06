# Kubernetes Aggregated API Server CA (resources/generated/tls/aggregator-ca.crt)
#
# TODO(diegs): this should be a sibling of the `--tls-ca-file` CA. However the
# self-signed CA ca.crt serves double-duty as the parent CA of other certs and
# as the `--tls-ca-file` CA. In the future that should be a separate CA as well.
resource "tls_private_key" "aggregator_ca" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "aggregator_ca" {
  key_algorithm   = "${tls_private_key.aggregator_ca.algorithm}"
  private_key_pem = "${tls_private_key.aggregator_ca.private_key_pem}"

  subject {
    common_name  = "aggregator"
    organization = "bootkube"
  }
}

resource "tls_locally_signed_cert" "aggregator_ca" {
  cert_request_pem = "${tls_cert_request.aggregator_ca.cert_request_pem}"

  is_ca_certificate     = true
  ca_key_algorithm      = "${tls_self_signed_cert.kube_ca.key_algorithm}"
  ca_private_key_pem    = "${tls_private_key.kube_ca.private_key_pem}"
  ca_cert_pem           = "${tls_self_signed_cert.kube_ca.cert_pem}"
  validity_period_hours = "${var.validity_period}"

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

resource "local_file" "aggregator_ca_key" {
  content  = "${tls_private_key.aggregator_ca.private_key_pem}"
  filename = "./generated/tls/aggregator-ca.key"
}

resource "local_file" "aggregator_ca_crt" {
  content  = "${tls_locally_signed_cert.aggregator_ca.cert_pem}"
  filename = "./generated/tls/aggregator-ca.crt"
}
