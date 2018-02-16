resource "local_file" "root_ca_cert" {
  content  = "${file(var.root_ca_cert_pem_path)}"
  filename = "./generated/tls/root-ca.crt"
}

resource "local_file" "kube_ca_key" {
  content  = "${file(var.kube_ca_key_pem_path)}"
  filename = "./generated/tls/kube-ca.key"
}

resource "local_file" "kube_ca_cert" {
  content  = "${file(var.kube_ca_cert_pem_path)}"
  filename = "./generated/tls/kube-ca.crt"
}

resource "local_file" "aggregator_ca_key" {
  content  = "${file(var.aggregator_key_pem_path)}"
  filename = "./generated/tls/aggregator-ca.key"
}

resource "local_file" "aggregator_ca_cert" {
  content  = "${file(var.aggregator_cert_pem_path)}"
  filename = "./generated/tls/aggregator-ca.crt"
}

resource "local_file" "etcd_ca_cert" {
  content  = "${file(var.etcd_ca_cert_pem_path)}"
  filename = "./generated/tls/etcd-client-ca.crt"
}
