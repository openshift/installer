output "root_ca_cert_pem" {
  value = "${var.root_ca_cert_pem == "" ? join("", tls_self_signed_cert.root_ca.*.cert_pem) : var.root_ca_cert_pem}"
}

output "kube_ca_cert_pem" {
  value = "${tls_locally_signed_cert.kube_ca.cert_pem}"
}

output "kube_ca_key_alg" {
  value = "${tls_locally_signed_cert.kube_ca.ca_key_algorithm}"
}

output "kube_ca_key_pem" {
  value = "${tls_private_key.kube_ca.private_key_pem}"
}

output "aggregator_ca_cert_pem" {
  value = "${tls_locally_signed_cert.aggregator_ca.cert_pem}"
}

output "aggregator_ca_key_alg" {
  value = "${tls_locally_signed_cert.aggregator_ca.ca_key_algorithm}"
}

output "aggregator_ca_key_pem" {
  value = "${tls_private_key.aggregator_ca.private_key_pem}"
}

output "etcd_ca_cert_pem" {
  value = "${tls_locally_signed_cert.etcd_ca.cert_pem}"
}

output "etcd_ca_key_alg" {
  value = "${tls_locally_signed_cert.etcd_ca.ca_key_algorithm}"
}

output "etcd_ca_key_pem" {
  value = "${tls_private_key.etcd_ca.private_key_pem}"
}

output "id" {
  value = "${sha1("
  ${join(" ",
    list(local_file.root_ca_cert.id,
    local_file.kube_ca_key.id,
    local_file.kube_ca_cert.id,
    local_file.aggregator_ca_key.id,
    local_file.aggregator_ca_cert.id,
    local_file.etcd_ca_cert.id)
    )}
  ")}"
}
