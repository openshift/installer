output "aggregator_ca_cert_pem" {
  value = "${tls_locally_signed_cert.aggregator_ca.cert_pem}"
}

output "ca_cert_pem" {
  value = "${var.ca_cert_pem == "" ? element(concat(tls_self_signed_cert.kube_ca.*.cert_pem, list("")), 0) : var.ca_cert_pem}"
}

output "ca_key_alg" {
  value = "${var.ca_cert_pem == "" ? element(concat(tls_self_signed_cert.kube_ca.*.key_algorithm, list("")), 0) : var.ca_key_alg}"
}

output "ca_key_pem" {
  value = "${var.ca_cert_pem == "" ? element(concat(tls_private_key.kube_ca.*.private_key_pem, list("")), 0) : var.ca_key_pem}"
}

output "admin_cert_pem" {
  value = "${tls_locally_signed_cert.admin.cert_pem}"
}

output "admin_key_pem" {
  value = "${tls_private_key.admin.private_key_pem}"
}

output "apiserver_cert_pem" {
  value = "${tls_locally_signed_cert.apiserver.cert_pem}"
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
    list(local_file.aggregator_ca_key.id,
    local_file.aggregator_ca_crt.id,
    local_file.apiserver_key.id,
    local_file.apiserver_crt.id,
    local_file.apiserver_proxy_key.id,
    local_file.apiserver_proxy_crt.id,
    local_file.kube_ca_key.id,
    local_file.kube_ca_crt.id,
    local_file.admin_key.id,
    local_file.admin_crt.id,)
    )}
  ")}"
}
