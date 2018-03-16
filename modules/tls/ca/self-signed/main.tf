# Root CA (resources/generated/tls/{root-ca.crt})
resource "tls_private_key" "root_ca" {
  count = "${var.root_ca_cert_pem == "" ? 1 : 0}"

  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_self_signed_cert" "root_ca" {
  count = "${var.root_ca_cert_pem == "" ? 1 : 0}"

  key_algorithm   = "${tls_private_key.root_ca.algorithm}"
  private_key_pem = "${tls_private_key.root_ca.private_key_pem}"

  subject {
    common_name         = "root-ca"
    organization        = "${uuid()}"
    organizational_unit = "tectonic"
  }

  is_ca_certificate = true

  # root ca defaults to being valid for 10 years.
  validity_period_hours = "87600"

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]

  lifecycle {
    ignore_changes = ["subject"]
  }
}

# Intermediate etcd CA (resources/generated/tls/{etcd-ca.crt})
resource "tls_private_key" "etcd_ca" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "etcd_ca" {
  key_algorithm   = "${tls_private_key.etcd_ca.algorithm}"
  private_key_pem = "${tls_private_key.etcd_ca.private_key_pem}"

  subject {
    common_name         = "etcd-ca"
    organization        = "${uuid()}"
    organizational_unit = "etcd"
  }

  lifecycle {
    ignore_changes = ["subject"]
  }
}

resource "tls_locally_signed_cert" "etcd_ca" {
  cert_request_pem = "${tls_cert_request.etcd_ca.cert_request_pem}"

  ca_key_algorithm   = "${var.root_ca_cert_pem == "" ? join("", tls_self_signed_cert.root_ca.*.key_algorithm) : var.root_ca_key_alg}"
  ca_private_key_pem = "${var.root_ca_cert_pem == "" ? join("", tls_private_key.root_ca.*.private_key_pem) : var.root_ca_key_pem}"
  ca_cert_pem        = "${var.root_ca_cert_pem == "" ? join("", tls_self_signed_cert.root_ca.*.cert_pem) : var.root_ca_cert_pem}"
  is_ca_certificate  = true

  # intermediate certs are valid for 3 years.
  validity_period_hours = "26280"

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}

# Intermediate kube CA (resources/generated/tls/{kube-ca.crt,kube-ca.key})
resource "tls_private_key" "kube_ca" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "kube_ca" {
  key_algorithm   = "${tls_private_key.kube_ca.algorithm}"
  private_key_pem = "${tls_private_key.kube_ca.private_key_pem}"

  subject {
    common_name         = "kube-ca"
    organization        = "${uuid()}"
    organizational_unit = "bootkube"
  }

  lifecycle {
    ignore_changes = ["subject"]
  }
}

resource "tls_locally_signed_cert" "kube_ca" {
  cert_request_pem = "${tls_cert_request.kube_ca.cert_request_pem}"

  ca_key_algorithm   = "${var.root_ca_cert_pem == "" ? join("", tls_self_signed_cert.root_ca.*.key_algorithm) : var.root_ca_key_alg}"
  ca_private_key_pem = "${var.root_ca_cert_pem == "" ? join("", tls_private_key.root_ca.*.private_key_pem) : var.root_ca_key_pem}"
  ca_cert_pem        = "${var.root_ca_cert_pem == "" ? join("", tls_self_signed_cert.root_ca.*.cert_pem) : var.root_ca_cert_pem}"
  is_ca_certificate  = true

  # intermediate certs are valid for 3 years.
  validity_period_hours = "26280"

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}

# Intermediate aggregator CA (resources/generated/tls/{aggregator-ca.crt,aggregator-ca.key})
resource "tls_private_key" "aggregator_ca" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "aggregator_ca" {
  key_algorithm   = "${tls_private_key.aggregator_ca.algorithm}"
  private_key_pem = "${tls_private_key.aggregator_ca.private_key_pem}"

  subject {
    common_name         = "aggregator"
    organization        = "${uuid()}"
    organizational_unit = "bootkube"
  }

  lifecycle {
    ignore_changes = ["subject"]
  }
}

resource "tls_locally_signed_cert" "aggregator_ca" {
  cert_request_pem = "${tls_cert_request.aggregator_ca.cert_request_pem}"

  ca_key_algorithm   = "${var.root_ca_cert_pem == "" ? join("", tls_self_signed_cert.root_ca.*.key_algorithm) : var.root_ca_key_alg}"
  ca_private_key_pem = "${var.root_ca_cert_pem == "" ? join("", tls_private_key.root_ca.*.private_key_pem) : var.root_ca_key_pem}"
  ca_cert_pem        = "${var.root_ca_cert_pem == "" ? join("", tls_self_signed_cert.root_ca.*.cert_pem) : var.root_ca_cert_pem}"
  is_ca_certificate  = true

  # intermediate certs are valid for 3 years.
  validity_period_hours = "26280"

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}
