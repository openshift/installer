# etcd assets
data "template_file" "etcd_ca_cert_pem" {
  template = "${tls_self_signed_cert.etcd_ca.cert_pem}"
}

data "template_file" "etcd_client_crt" {
  template = "${tls_locally_signed_cert.etcd_client.cert_pem}"
}

data "template_file" "etcd_client_key" {
  template = "${tls_private_key.etcd_client.private_key_pem}"
}

resource "local_file" "etcd_ca_crt" {
  content  = "${data.template_file.etcd_ca_cert_pem.rendered}"
  filename = "./generated/tls/etcd-ca.crt"
}

resource "local_file" "etcd_client_crt" {
  content  = "${data.template_file.etcd_client_crt.rendered}"
  filename = "./generated/tls/etcd-client.crt"
}

resource "local_file" "etcd_client_key" {
  content  = "${data.template_file.etcd_client_key.rendered}"
  filename = "./generated/tls/etcd-client.key"
}

resource "local_file" "etcd_server_crt" {
  content  = "${tls_locally_signed_cert.etcd_server.cert_pem}"
  filename = "./generated/tls/etcd-server.crt"
}

resource "local_file" "etcd_server_key" {
  content  = "${tls_private_key.etcd_server.private_key_pem}"
  filename = "./generated/tls/etcd-server.key"
}

resource "local_file" "etcd_peer_crt" {
  content  = "${tls_locally_signed_cert.etcd_peer.cert_pem}"
  filename = "./generated/tls/etcd-peer.crt"
}

resource "local_file" "etcd_peer_key" {
  content  = "${tls_private_key.etcd_peer.private_key_pem}"
  filename = "./generated/tls/etcd-peer.key"
}
