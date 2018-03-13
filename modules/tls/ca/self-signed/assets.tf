resource "local_file" "root_ca_cert" {
  content  = "${var.root_ca_cert_pem == "" ? join("", tls_self_signed_cert.root_ca.*.cert_pem) : var.root_ca_cert_pem}"
  filename = "./generated/tls/root-ca.crt"
}

resource "local_file" "kube_ca_key" {
  content  = "${tls_private_key.kube_ca.private_key_pem}"
  filename = "./generated/tls/kube-ca.key"
}

resource "local_file" "kube_ca_cert" {
  content  = "${tls_locally_signed_cert.kube_ca.cert_pem}"
  filename = "./generated/tls/kube-ca.crt"
}

resource "local_file" "aggregator_ca_key" {
  content  = "${tls_private_key.aggregator_ca.private_key_pem}"
  filename = "./generated/tls/aggregator-ca.key"
}

resource "local_file" "aggregator_ca_cert" {
  content  = "${tls_locally_signed_cert.aggregator_ca.cert_pem}"
  filename = "./generated/tls/aggregator-ca.crt"
}

resource "local_file" "etcd_ca_cert" {
  content  = "${tls_locally_signed_cert.etcd_ca.cert_pem}"
  filename = "./generated/tls/etcd-client-ca.crt"
}
