# NOTE: Across this module, the following syntax is used at various places:
#   `"${var.ca_cert == "" ? join(" ", tls_private_key.kube-ca.*.private_key_pem) : var.ca_key}"`
#
# Due to https://github.com/hashicorp/hil/issues/50, both sides of conditions
# are evaluated, until one of them is discarded. Unfortunately, the
# `{tls_private_key/tls_self_signed_cert}.kube-ca` resources are created
# conditionally and might not be present - in which case an error is
# generated. Because a `count` is used on these ressources, the resources can be
# referenced as lists with the `.*` notation, and arrays are allowed to be
# empty. The `join()` interpolation function is then used to cast them back to
# a string. Since `count` can only be 0 or 1, the returned value is either empty
# (and discarded anyways) or the desired value.

# Kubernetes CA (resources/generated/tls/{ca.crt,ca.key})
resource "tls_private_key" "kube_ca" {
  count = "${var.ca_cert == "" ? 1 : 0}"

  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_self_signed_cert" "kube_ca" {
  count = "${var.ca_cert == "" ? 1 : 0}"

  key_algorithm   = "${tls_private_key.kube_ca.algorithm}"
  private_key_pem = "${tls_private_key.kube_ca.private_key_pem}"

  subject {
    common_name  = "kube-ca"
    organization = "bootkube"
  }

  is_ca_certificate     = true
  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}

resource "local_file" "kube_ca_key" {
  content  = "${var.ca_cert == "" ? join(" ", tls_private_key.kube_ca.*.private_key_pem) : var.ca_key}"
  filename = "./generated/tls/ca.key"
}

resource "local_file" "kube_ca_crt" {
  content  = "${var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube_ca.*.cert_pem) : var.ca_cert}"
  filename = "./generated/tls/ca.crt"
}

# Kubernetes API Server (resources/generated/tls/{apiserver.key,apiserver.crt})
resource "tls_private_key" "apiserver" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "apiserver" {
  key_algorithm   = "${tls_private_key.apiserver.algorithm}"
  private_key_pem = "${tls_private_key.apiserver.private_key_pem}"

  subject {
    common_name  = "kube-apiserver"
    organization = "kube-master"
  }

  dns_names = [
    "${replace(element(split(":", var.kube_apiserver_url), 1), "/", "")}",
    "kubernetes",
    "kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.cluster.local",
  ]

  ip_addresses = [
    "${cidrhost(var.service_cidr, 1)}",
  ]
}

resource "tls_locally_signed_cert" "apiserver" {
  cert_request_pem = "${tls_cert_request.apiserver.cert_request_pem}"

  ca_key_algorithm   = "${var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube_ca.*.key_algorithm) : var.ca_key_alg}"
  ca_private_key_pem = "${var.ca_cert == "" ? join(" ", tls_private_key.kube_ca.*.private_key_pem) : var.ca_key}"
  ca_cert_pem        = "${var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube_ca.*.cert_pem): var.ca_cert}"

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

resource "local_file" "apiserver_key" {
  content  = "${tls_private_key.apiserver.private_key_pem}"
  filename = "./generated/tls/apiserver.key"
}

resource "local_file" "apiserver_crt" {
  content  = "${tls_locally_signed_cert.apiserver.cert_pem}"
  filename = "./generated/tls/apiserver.crt"
}

# Kubernete's Service Account (resources/generated/tls/{service-account.key,service-account.pub})
resource "tls_private_key" "service_account" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "local_file" "service_account_key" {
  content  = "${tls_private_key.service_account.private_key_pem}"
  filename = "./generated/tls/service-account.key"
}

resource "local_file" "service_account_crt" {
  content  = "${tls_private_key.service_account.public_key_pem}"
  filename = "./generated/tls/service-account.pub"
}

# Kubelet
resource "tls_private_key" "kubelet" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "kubelet" {
  key_algorithm   = "${tls_private_key.kubelet.algorithm}"
  private_key_pem = "${tls_private_key.kubelet.private_key_pem}"

  subject {
    common_name  = "kubelet"
    organization = "system:masters"
  }
}

resource "tls_locally_signed_cert" "kubelet" {
  cert_request_pem = "${tls_cert_request.kubelet.cert_request_pem}"

  ca_key_algorithm   = "${var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube_ca.*.key_algorithm) : var.ca_key_alg}"
  ca_private_key_pem = "${var.ca_cert == "" ? join(" ", tls_private_key.kube_ca.*.private_key_pem) : var.ca_key}"
  ca_cert_pem        = "${var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube_ca.*.cert_pem) : var.ca_cert}"

  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

resource "local_file" "kubelet_key" {
  content  = "${tls_private_key.kubelet.private_key_pem}"
  filename = "./generated/tls/kubelet.key"
}

resource "local_file" "kubelet_crt" {
  content  = "${tls_locally_signed_cert.kubelet.cert_pem}"
  filename = "./generated/tls/kubelet.crt"
}
