output "admin_cert_pem" {
  value = "${tls_locally_signed_cert.admin.cert_pem}"
}

output "admin_key_pem" {
  value = "${tls_private_key.admin.private_key_pem}"
}

output "apiserver_cert_pem" {
  value = "${data.template_file.apiserver_cert.rendered}"
}

output "apiserver_key_pem" {
  value = "${tls_private_key.apiserver.private_key_pem}"
}

output "apiserver_proxy_cert_pem" {
  value = "${tls_locally_signed_cert.apiserver_proxy.cert_pem}"
}

output "apiserver_proxy_key_pem" {
  value = "${tls_private_key.apiserver_proxy.private_key_pem}"
}

output "id" {
  value = "${sha1("
  ${join(" ",
    list(local_file.apiserver_key.id,
    local_file.apiserver_cert.id,
    local_file.apiserver_proxy_key.id,
    local_file.apiserver_proxy_cert.id,
    local_file.admin_key.id,
    local_file.admin_cert.id,)
    )}
  ")}"
}
