output "etcd_client_cert_pem" {
  value = "${tls_locally_signed_cert.etcd_client.cert_pem}"
}

output "etcd_client_key_pem" {
  value = "${tls_private_key.etcd_client.private_key_pem}"
}

output "etcd_peer_cert_pem" {
  value = "${tls_locally_signed_cert.etcd_peer.cert_pem}"
}

output "etcd_peer_key_pem" {
  value = "${tls_private_key.etcd_peer.private_key_pem}"
}

output "etcd_server_cert_pem" {
  value = "${tls_locally_signed_cert.etcd_server.cert_pem}"
}

output "etcd_server_key_pem" {
  value = "${tls_private_key.etcd_server.private_key_pem}"
}

output "id" {
  value = "${sha1("
  ${join(" ",
    list(local_file.etcd_server_cert.id,
    local_file.etcd_server_key.id,
    local_file.etcd_client_cert.id,
    local_file.etcd_client_key.id,
    local_file.etcd_peer_cert.id,
    local_file.etcd_peer_key.id)
    )}
  ")}"
}
