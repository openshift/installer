output "etcd_client_cert_pem" {
  value = "${tls_locally_signed_cert.etcd_client.cert_pem}"
}

output "etcd_client_key_pem" {
  value = "${tls_private_key.etcd_client.private_key_pem}"
}

output "id" {
  value = "${sha1("
  ${join(" ",
    list(local_file.etcd_client_cert.id,
    local_file.etcd_client_key.id)
    )}
  ")}"
}
