output "root_ca_cert_pem" {
  value = "${var.root_ca_cert_pem_path == "" ? join("", tls_self_signed_cert.root_ca.*.cert_pem) : file(local._root_ca_cert_pem_path)}"
}

output "kube_ca_cert_pem" {
  value = "${var.kube_ca_cert_pem_path == "" ? join("", tls_locally_signed_cert.kube_ca.*.cert_pem) : file(local._kube_ca_cert_pem_path)}"
}

output "kube_ca_key_alg" {
  value = "${var.kube_ca_key_alg == "" ? join("", tls_locally_signed_cert.kube_ca.*.ca_key_algorithm) : var.kube_ca_key_alg}"
}

output "kube_ca_key_pem" {
  value = "${var.kube_ca_key_pem_path == "" ? join("", tls_private_key.kube_ca.*.private_key_pem) : file(local._kube_ca_key_pem_path)}"
}

output "aggregator_ca_cert_pem" {
  value = "${var.aggregator_ca_cert_pem_path == "" ? join("", tls_locally_signed_cert.aggregator_ca.*.cert_pem) : file(local._aggregator_ca_cert_pem_path)}"
}

output "aggregator_ca_key_alg" {
  value = "${var.aggregator_ca_key_alg == "" ? join("", tls_locally_signed_cert.aggregator_ca.*.ca_key_algorithm) : var.aggregator_ca_key_alg}"
}

output "aggregator_ca_key_pem" {
  value = "${var.aggregator_ca_key_pem_path == "" ? join("", tls_private_key.aggregator_ca.*.private_key_pem) : file(local._aggregator_ca_key_pem_path)}"
}

output "service_serving_ca_cert_pem" {
  value = "${var.service_serving_ca_cert_pem_path == "" ? join("", tls_locally_signed_cert.service_serving_ca.*.cert_pem) : file(local._service_serving_ca_cert_pem_path)}"
}

output "service_serving_ca_key_alg" {
  value = "${var.service_serving_ca_key_alg == "" ? join("", tls_locally_signed_cert.service_serving_ca.*.ca_key_algorithm) : var.service_serving_ca_key_alg}"
}

output "service_serving_ca_key_pem" {
  value = "${var.service_serving_ca_key_pem_path == "" ? join("", tls_private_key.service_serving_ca.*.private_key_pem) : file(local._service_serving_ca_key_pem_path)}"
}

output "etcd_ca_cert_pem" {
  value = "${var.etcd_ca_cert_pem_path == "" ? join("", tls_locally_signed_cert.etcd_ca.*.cert_pem) : file(local._etcd_ca_cert_pem_path)}"
}

output "etcd_ca_key_alg" {
  value = "${var.etcd_ca_key_alg == "" ? join("", tls_locally_signed_cert.etcd_ca.*.ca_key_algorithm) : var.etcd_ca_key_alg}"
}

output "etcd_ca_key_pem" {
  value = "${var.etcd_ca_key_pem_path == "" ? join("", tls_private_key.etcd_ca.*.private_key_pem) : file(local._etcd_ca_key_pem_path)}"
}

output "id" {
  value = "${sha1("
  ${join(" ",
    list(local_file.root_ca_cert.id,
    local_file.kube_ca_key.id,
    local_file.kube_ca_cert.id,
    local_file.aggregator_ca_key.id,
    local_file.aggregator_ca_cert.id,
    local_file.etcd_ca_key.id,
    local_file.etcd_ca_cert.id)
    )}
  ")}"
}
