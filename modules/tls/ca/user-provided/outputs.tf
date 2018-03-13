output "root_ca_cert_pem" {
  value = "${file(var.root_ca_cert_pem_path)}"
}

output "root_ca_key_pem" {
  value = "${file(var.kube_ca_key_pem_path)}"
}

output "kube_ca_cert_pem" {
  value = "${file(var.kube_ca_cert_pem_path)}"
}

output "kube_ca_key_pem" {
  value = "${file(var.aggregator_key_pem_path)}"
}

output "aggregator_ca_cert_pem" {
  value = "${file(var.aggregator_cert_pem_path)}"
}

output "aggregator_ca_key_pem" {
  value = "${file(var.aggregator_key_pem_path)}"
}

output "etcd_ca_cert_pem" {
  value = "${file(var.etcd_ca_cert_pem_path)}"
}
