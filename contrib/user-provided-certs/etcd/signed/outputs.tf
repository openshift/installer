output "etcd_ca_crt_pem" {
  value = "${data.template_file.etcd_ca_cert_pem.rendered}"
}

output "etcd_client_crt_pem" {
  value = "${data.template_file.etcd_client_crt.rendered}"
}

output "etcd_client_key_pem" {
  value = "${data.template_file.etcd_client_key.rendered}"
}

output "etcd_peer_crt_pem" {
  value = "${tls_locally_signed_cert.etcd_peer.cert_pem}"
}

output "etcd_peer_key_pem" {
  value = "${tls_private_key.etcd_peer.private_key_pem}"
}

output "etcd_server_crt_pem" {
  value = "${tls_locally_signed_cert.etcd_server.cert_pem}"
}

output "etcd_server_key_pem" {
  value = "${tls_private_key.etcd_server.private_key_pem}"
}
