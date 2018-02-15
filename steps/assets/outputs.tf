output "kube_dns_service_ip" {
  value = "${module.bootkube.kube_dns_service_ip}"
}

output "kubeconfig_content" {
  value = "${module.bootkube.kubeconfig}"
}

output "kubeconfig_kubelet_content" {
  value = "${module.bootkube.kubeconfig-kubelet}"
}

output "cluster_id" {
  value = "${module.tectonic.cluster_id}"
}

# TLS
output "etcd_ca_crt_pem" {
  value = "${module.etcd_certs.etcd_ca_crt_pem}"
}

output "etcd_client_crt_pem" {
  value = "${module.etcd_certs.etcd_client_crt_pem}"
}

output "etcd_client_key_pem" {
  value = "${module.etcd_certs.etcd_client_key_pem}"
}

output "etcd_peer_crt_pem" {
  value = "${module.etcd_certs.etcd_peer_crt_pem}"
}

output "etcd_peer_key_pem" {
  value = "${module.etcd_certs.etcd_peer_key_pem}"
}

output "etcd_server_crt_pem" {
  value = "${module.etcd_certs.etcd_server_crt_pem}"
}

output "etcd_server_key_pem" {
  value = "${module.etcd_certs.etcd_server_key_pem}"
}

output "ingress_certs_ca_cert_pem" {
  value = "${module.ingress_certs.ca_cert_pem}"
}

output "kube_certs_ca_cert_pem" {
  value = "${module.kube_certs.ca_cert_pem}"
}
