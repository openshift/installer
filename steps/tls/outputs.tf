output "admin_cert_pem" {
  value = "${module.kube_certs.admin_cert_pem}"
}

output "admin_key_pem" {
  value = "${module.kube_certs.admin_key_pem}"
}

output "aggregator_ca_cert_pem" {
  value = "${module.ca_certs.aggregator_ca_cert_pem}"
}

output "aggregator_ca_key_pem" {
  value = "${module.ca_certs.aggregator_ca_key_pem}"
}

output "apiserver_cert_pem" {
  value = "${module.kube_certs.apiserver_cert_pem}"
}

output "apiserver_key_pem" {
  value = "${module.kube_certs.apiserver_key_pem}"
}

output "apiserver_proxy_cert_pem" {
  value = "${module.kube_certs.apiserver_proxy_cert_pem}"
}

output "etcd_ca_cert_pem" {
  value = "${module.ca_certs.etcd_ca_cert_pem}"
}

output "etcd_ca_key_pem" {
  value = "${module.ca_certs.etcd_ca_key_pem}"
}

output "apiserver_proxy_key_pem" {
  value = "${module.kube_certs.apiserver_proxy_key_pem}"
}

output "kubelet_cert_pem" {
  value = "${module.kube_certs.kubelet_cert_pem}"
}

output "kubelet_key_pem" {
  value = "${module.kube_certs.kubelet_key_pem}"
}

output "etcd_client_cert_pem" {
  value = "${module.etcd_certs.etcd_client_cert_pem}"
}

output "etcd_client_key_pem" {
  value = "${module.etcd_certs.etcd_client_key_pem}"
}

output "identity_client_ca_cert" {
  value = "${module.ca_certs.root_ca_cert_pem}"
}

output "identity_client_cert_pem" {
  value = "${module.identity_certs.client_cert_pem}"
}

output "identity_client_key_pem" {
  value = "${module.identity_certs.client_key_pem}"
}

output "identity_server_ca_cert" {
  value = "${module.ca_certs.kube_ca_cert_pem}"
}

output "identity_server_cert_pem" {
  value = "${module.identity_certs.server_cert_pem}"
}

output "identity_server_key_pem" {
  value = "${module.identity_certs.server_key_pem}"
}

output "ingress_ca_cert_pem" {
  value = "${module.ingress_certs.ca_cert_pem}"
}

output "ingress_cert_pem" {
  value = "${module.ingress_certs.cert_pem}"
}

output "ingress_key_pem" {
  value = "${module.ingress_certs.key_pem}"
}

output "kube_ca_cert_pem" {
  value = "${module.ca_certs.kube_ca_cert_pem}"
}

output "kube_ca_key_pem" {
  value = "${module.ca_certs.kube_ca_key_pem}"
}

output "oidc_ca_cert" {
  value = "${module.ingress_certs.ca_cert_pem}"
}

output "root_ca_cert_pem" {
  value = "${module.ca_certs.root_ca_cert_pem}"
}
