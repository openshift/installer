output "admin_cert_pem" {
  value = "${tls_locally_signed_cert.admin.cert_pem}"
}

output "admin_key_pem" {
  value = "${tls_private_key.admin.private_key_pem}"
}

output "kubelet_cert_pem" {
  value = "${tls_locally_signed_cert.kubelet.cert_pem}"
}

output "kubelet_key_pem" {
  value = "${tls_private_key.kubelet.private_key_pem}"
}

output "apiserver_cert_pem" {
  value = "${data.template_file.apiserver_cert.rendered}"
}

output "apiserver_key_pem" {
  value = "${tls_private_key.apiserver.private_key_pem}"
}

output "openshift_apiserver_cert_pem" {
  value = "${data.template_file.openshift_apiserver_cert.rendered}"
}

output "openshift_apiserver_key_pem" {
  value = "${tls_private_key.openshift_apiserver.private_key_pem}"
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
    local_file.openshift_apiserver_key.id,
    local_file.openshift_apiserver_cert.id,
    local_file.apiserver_proxy_key.id,
    local_file.apiserver_proxy_cert.id,
    local_file.admin_key.id,
    local_file.admin_cert.id,
    local_file.kubelet_key.id,
    local_file.kubelet_cert.id,)
    )}
  ")}"
}
