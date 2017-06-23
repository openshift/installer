# Cryptographically-secure ramdon strings used by various components.

resource "random_id" "admin_user_id" {
  byte_length = 16
}

resource "random_id" "kubectl_secret" {
  byte_length = 16
}

resource "random_id" "console_secret" {
  byte_length = 16
}

resource "random_id" "tectonic_monitoring_auth_cookie_secret" {
  byte_length = 16
}

# Ingress' server certificate

resource "tls_private_key" "ingress" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "ingress" {
  key_algorithm   = "${tls_private_key.ingress.algorithm}"
  private_key_pem = "${tls_private_key.ingress.private_key_pem}"

  subject {
    common_name = "${element(split(":", var.base_address), 0)}"
  }

  # subject commonName is deprecated per RFC2818 in favor of
  # subjectAltName
  dns_names = [
    "${element(split(":", var.base_address), 0)}",
  ]
}

resource "tls_locally_signed_cert" "ingress" {
  cert_request_pem = "${tls_cert_request.ingress.cert_request_pem}"

  ca_key_algorithm   = "${var.ca_key_alg}"
  ca_private_key_pem = "${var.ca_key}"
  ca_cert_pem        = "${var.ca_cert}"

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

# Identity's gRPC server/client certificates

resource "tls_private_key" "identity-server" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "identity-server" {
  key_algorithm   = "${tls_private_key.identity-server.algorithm}"
  private_key_pem = "${tls_private_key.identity-server.private_key_pem}"

  subject {
    common_name = "tectonic-identity-api.tectonic-system.svc.cluster.local"
  }
}

resource "tls_locally_signed_cert" "identity-server" {
  cert_request_pem = "${tls_cert_request.identity-server.cert_request_pem}"

  ca_key_algorithm   = "${var.ca_key_alg}"
  ca_private_key_pem = "${var.ca_key}"
  ca_cert_pem        = "${var.ca_cert}"

  validity_period_hours = 8760

  allowed_uses = [
    "server_auth",
  ]
}

resource "tls_private_key" "identity-client" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "identity-client" {
  key_algorithm   = "${tls_private_key.identity-client.algorithm}"
  private_key_pem = "${tls_private_key.identity-client.private_key_pem}"

  subject {
    common_name = "tectonic-identity-api.tectonic-system.svc.cluster.local"
  }
}

resource "tls_locally_signed_cert" "identity-client" {
  cert_request_pem = "${tls_cert_request.identity-client.cert_request_pem}"

  ca_key_algorithm   = "${var.ca_key_alg}"
  ca_private_key_pem = "${var.ca_key}"
  ca_cert_pem        = "${var.ca_cert}"

  validity_period_hours = 8760

  allowed_uses = [
    "client_auth",
  ]
}
