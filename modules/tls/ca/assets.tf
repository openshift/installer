resource "local_file" "root_ca_cert" {
  content  = "${var.root_ca_cert_pem_path == "" ? join("", tls_self_signed_cert.root_ca.*.cert_pem) : file(local._root_ca_cert_pem_path )}"
  filename = "./generated/tls/root-ca.crt"
}

resource "local_file" "kube_ca_key" {
  content  = "${var.kube_ca_key_pem_path == "" ? join("", tls_private_key.kube_ca.*.private_key_pem) : file(local._kube_ca_key_pem_path)}"
  filename = "./generated/tls/kube-ca.key"
}

resource "local_file" "kube_ca_cert" {
  content  = "${var.kube_ca_cert_pem_path == "" ? join("", tls_locally_signed_cert.kube_ca.*.cert_pem) : file(local._kube_ca_cert_pem_path )}"
  filename = "./generated/tls/kube-ca.crt"
}

resource "local_file" "aggregator_ca_key" {
  content  = "${var.aggregator_ca_key_pem_path == "" ? join("", tls_private_key.aggregator_ca.*.private_key_pem) : file(local._aggregator_ca_key_pem_path)}"
  filename = "./generated/tls/aggregator-ca.key"
}

resource "local_file" "aggregator_ca_cert" {
  content  = "${var.aggregator_ca_cert_pem_path == "" ? join("", tls_locally_signed_cert.aggregator_ca.*.cert_pem) : file(local._aggregator_ca_cert_pem_path)}"
  filename = "./generated/tls/aggregator-ca.crt"
}

resource "local_file" "service_serving_ca_key" {
  content  = "${var.service_serving_ca_key_pem_path == "" ? join("", tls_private_key.service_serving_ca.*.private_key_pem) : file(local._service_serving_ca_key_pem_path)}"
  filename = "./generated/tls/service-serving-ca.key"
}

resource "local_file" "service_serving_ca_cert" {
  content  = "${var.service_serving_ca_cert_pem_path == "" ? join("", tls_locally_signed_cert.service_serving_ca.*.cert_pem) : file(local._service_serving_ca_cert_pem_path)}"
  filename = "./generated/tls/service-serving-ca.crt"
}

resource "local_file" "etcd_ca_key" {
  content  = "${var.etcd_ca_key_pem_path == "" ? join("", tls_private_key.etcd_ca.*.private_key_pem) : file(local._etcd_ca_key_pem_path)}"
  filename = "./generated/tls/etcd-client-ca.key"
}

resource "local_file" "etcd_ca_cert" {
  content  = "${var.etcd_ca_cert_pem_path == "" ? join("", tls_locally_signed_cert.etcd_ca.*.cert_pem) : file(local._etcd_ca_cert_pem_path)}"
  filename = "./generated/tls/etcd-client-ca.crt"
}
