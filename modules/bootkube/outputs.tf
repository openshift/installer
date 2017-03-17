output "kubeconfig" {
  value = "${data.template_file.kubeconfig.rendered}"
}

output "ca_cert" {
  value = "${var.ca_cert == "" ? tls_self_signed_cert.kube-ca.cert_pem : var.ca_cert}"
}

output "ca_key_alg" {
  value = "${var.ca_cert == "" ? tls_self_signed_cert.kube-ca.key_algorithm : var.ca_key_alg}"
}

output "ca_key" {
  value = "${var.ca_cert == "" ? tls_private_key.kube-ca.private_key_pem : var.ca_key}"
}
