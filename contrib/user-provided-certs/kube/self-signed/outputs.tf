output "aggregator_ca_cert_pem" {
  value = "${tls_self_signed_cert.aggregator_ca.cert_pem}"
}

output "ca_cert_pem" {
  value = "${tls_self_signed_cert.kube_ca.cert_pem}"
}

output "ca_key_alg" {
  value = "${tls_self_signed_cert.kube_ca.key_algorithm}"
}

output "ca_key_pem" {
  value = "${tls_private_key.kube_ca.private_key_pem}"
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
