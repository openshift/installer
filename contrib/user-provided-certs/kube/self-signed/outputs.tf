output "ca_cert_pem" {
  value = "${tls_self_signed_cert.kube_ca.cert_pem}"
}

output "ca_key_alg" {
  value = "${tls_self_signed_cert.kube_ca.key_algorithm}"
}

output "ca_key_pem" {
  value = "${tls_private_key.kube_ca.private_key_pem}"
}

output "kubelet_cert_pem" {
  value = "${tls_locally_signed_cert.kubelet.cert_pem}"
}

output "kubelet_key_pem" {
  value = "${tls_private_key.kubelet.private_key_pem}"
}

output "apiserver_cert_pem" {
  value = "${tls_locally_signed_cert.apiserver.cert_pem}"
}

output "apiserver_key_pem" {
  value = "${tls_private_key.apiserver.private_key_pem}"
}
