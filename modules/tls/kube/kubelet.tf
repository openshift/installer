# kubelet (generated/tls/{kubelet.key,kubelet.crt})
# Used to create kubeconfig (generated/auth/kubeconfig-kubelet) with CSR only privileges.
resource "tls_private_key" "kubelet" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

resource "tls_cert_request" "kubelet" {
  key_algorithm   = "${tls_private_key.kubelet.algorithm}"
  private_key_pem = "${tls_private_key.kubelet.private_key_pem}"

  subject {
    common_name  = "system:serviceaccount:kube-system:default"
    organization = "system:serviceaccounts:kube-system"
  }
}

resource "tls_locally_signed_cert" "kubelet" {
  cert_request_pem = "${tls_cert_request.kubelet.cert_request_pem}"

  ca_key_algorithm   = "${var.kube_ca_key_alg}"
  ca_private_key_pem = "${var.kube_ca_key_pem}"
  ca_cert_pem        = "${var.kube_ca_cert_pem}"

  # want bootstrap node to rotate certificate as soon as possible for the first time
  validity_period_hours = "1"

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}
