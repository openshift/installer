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

resource "tls_locally_signed_cert" "apiserver_proxy" {
  cert_request_pem = "${tls_cert_request.apiserver_proxy.cert_request_pem}"

  ca_key_algorithm   = "${tls_locally_signed_cert.aggregator_ca.ca_key_algorithm}"
  ca_private_key_pem = "${tls_private_key.aggregator_ca.private_key_pem}"
  ca_cert_pem        = "${tls_locally_signed_cert.aggregator_ca.cert_pem}"

  validity_period_hours = 26280

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

resource "local_file" "apiserver_proxy_key" {
  content  = "${tls_private_key.apiserver_proxy.private_key_pem}"
  filename = "./generated/tls/apiserver-proxy.key"
}

resource "local_file" "apiserver_proxy_crt" {
  content  = "${tls_locally_signed_cert.apiserver_proxy.cert_pem}"
  filename = "./generated/tls/apiserver-proxy.crt"
}
