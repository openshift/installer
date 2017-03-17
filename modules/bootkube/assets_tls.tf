# Kubernetes CA (resources/generated/tls/{ca.crt,ca.key})
resource "tls_private_key" "kube-ca" {
  count = "${var.ca_cert == "" ? 1 : 0}"

  algorithm = "RSA"
  rsa_bits = "2048"
}

resource "tls_self_signed_cert" "kube-ca" {
  count = "${var.ca_cert == "" ? 1 : 0}"

  key_algorithm = "${tls_private_key.kube-ca.algorithm}"
  private_key_pem = "${tls_private_key.kube-ca.private_key_pem}"

  subject {
    common_name = "kube-ca"
    organization = "bootkube"
  }

  is_ca_certificate = true
  validity_period_hours = 8760
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "cert_signing",
  ]
}

resource "localfile_file" "kube-ca-key" {
  content = "${var.ca_cert == "" ? tls_private_key.kube-ca.private_key_pem : var.ca_key}"
  destination = "${path.cwd}/generated/tls/ca.key"
}

resource "localfile_file" "kube-ca-crt" {
  content = "${var.ca_cert == "" ? tls_self_signed_cert.kube-ca.cert_pem : var.ca_cert}"
  destination = "${path.cwd}/generated/tls/ca.crt"
}

# Kubernetes API Server (resources/generated/tls/{apiserver.key,apiserver.crt})
resource "tls_private_key" "apiserver" {
  algorithm = "RSA"
  rsa_bits = "2048"
}

resource "tls_cert_request" "apiserver" {
  key_algorithm = "${tls_private_key.apiserver.algorithm}"
  private_key_pem = "${tls_private_key.apiserver.private_key_pem}"

  subject {
    common_name = "kube-apiserver"
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
    "${var.kube_apiserver_service_ip}",
  ]
}

resource "tls_locally_signed_cert" "apiserver" {
  cert_request_pem = "${tls_cert_request.apiserver.cert_request_pem}"

  ca_key_algorithm = "${var.ca_cert == "" ? tls_self_signed_cert.kube-ca.key_algorithm : var.ca_key_alg}"
  ca_private_key_pem = "${var.ca_cert == "" ? tls_private_key.kube-ca.private_key_pem : var.ca_key}"
  ca_cert_pem = "${var.ca_cert == "" ? tls_self_signed_cert.kube-ca.cert_pem : var.ca_cert}"

  validity_period_hours = 8760
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

resource "localfile_file" "apiserver-key" {
  content = "${tls_private_key.apiserver.private_key_pem}"
  destination = "${path.cwd}/generated/tls/apiserver.key"
}

resource "localfile_file" "apiserver-crt" {
  content = "${tls_locally_signed_cert.apiserver.cert_pem}"
  destination = "${path.cwd}/generated/tls/apiserver.crt"
}

# Kubernete's Service Account (resources/generated/tls/{service-account.key,service-account.pub})
resource "tls_private_key" "service-account" {
  algorithm = "RSA"
  rsa_bits = "2048"
}

resource "localfile_file" "service-account-key" {
  content = "${tls_private_key.service-account.private_key_pem}"
  destination = "${path.cwd}/generated/tls/service-account.key"
}

resource "localfile_file" "service-account-crt" {
  content = "${tls_private_key.service-account.public_key_pem}"
  destination = "${path.cwd}/generated/tls/service-account.pub"
}

# Kubelet
resource "tls_private_key" "kubelet" {
  algorithm = "RSA"
  rsa_bits = "2048"
}

resource "tls_cert_request" "kubelet" {
  key_algorithm = "${tls_private_key.kubelet.algorithm}"
  private_key_pem = "${tls_private_key.kubelet.private_key_pem}"

  subject {
    common_name = "kubelet"
    organization = "system:masters"
  }
}

resource "tls_locally_signed_cert" "kubelet" {
  cert_request_pem = "${tls_cert_request.kubelet.cert_request_pem}"

  ca_key_algorithm = "${var.ca_cert == "" ? tls_self_signed_cert.kube-ca.key_algorithm : var.ca_key_alg}"
  ca_private_key_pem = "${var.ca_cert == "" ? tls_private_key.kube-ca.private_key_pem : var.ca_key}"
  ca_cert_pem = "${var.ca_cert == "" ? tls_self_signed_cert.kube-ca.cert_pem : var.ca_cert}"

  validity_period_hours = 8760
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "client_auth",
  ]
}

resource "localfile_file" "kubelet-key" {
  content = "${tls_private_key.kubelet.private_key_pem}"
  destination = "${path.cwd}/generated/tls/kubelet.key"
}

resource "localfile_file" "kubelet-crt" {
  content = "${tls_locally_signed_cert.kubelet.cert_pem}"
  destination = "${path.cwd}/generated/tls/kubelet.crt"
}
