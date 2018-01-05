resource "tls_private_key" "kube_ca" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_self_signed_cert" "kube_ca" {
  key_algorithm   = "${tls_private_key.kube_ca.algorithm}"
  private_key_pem = "${tls_private_key.kube_ca.private_key_pem}"

  subject {
    common_name  = "kube-ca"
    organization = "bootkube"
  }

  is_ca_certificate     = true
  validity_period_hours = 26280

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}

resource "local_file" "kube_ca_key" {
  content  = "${tls_private_key.kube_ca.private_key_pem}"
  filename = "./generated/tls/ca.key"
}

resource "local_file" "kube_ca_crt" {
  content  = "${tls_self_signed_cert.kube_ca.cert_pem}"
  filename = "./generated/tls/ca.crt"
}
