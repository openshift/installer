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

# Openshift API Server (resources/generated/tls/{openshift-apiserver.key,openshift-apiserver.crt})
resource "tls_private_key" "openshift_apiserver" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "openshift_apiserver" {
  key_algorithm   = "${tls_private_key.openshift_apiserver.algorithm}"
  private_key_pem = "${tls_private_key.openshift_apiserver.private_key_pem}"

  subject {
    common_name  = "openshift-apiserver"
    organization = "kube-master"
  }

  dns_names = [
    "${replace(element(split(":", var.kube_apiserver_url), 1), "/", "")}",
    "openshift-apiserver",
    "openshift-apiserver.kube-system",
    "openshift-apiserver.kube-system.svc",
    "openshift-apiserver.kube-system.svc.cluster.local",
    "localhost",
    "127.0.0.1",
  ]

  ip_addresses = [
    "${cidrhost(var.service_cidr, 1)}",
  ]
}

resource "tls_locally_signed_cert" "openshift_apiserver" {
  cert_request_pem = "${tls_cert_request.openshift_apiserver.cert_request_pem}"

  ca_key_algorithm      = "${var.aggregator_ca_key_alg}"
  ca_private_key_pem    = "${var.aggregator_ca_key_pem}"
  ca_cert_pem           = "${var.aggregator_ca_cert_pem}"
  validity_period_hours = "26280"

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}
