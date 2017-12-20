module "kube_certs" {
  source = "./kube/self-signed"

  kube_apiserver_url = "https://${var.api_fqdn}:443"
  service_cidr       = "${var.service_cidr}"
}

module "etcd_certs" {
  source = "./etcd/signed"

  etcd_cert_dns_names = "${var.etcd_dns_names}"
  service_cidr        = "${var.service_cidr}"
}

module "ingress_certs" {
  source = "./ingress/self-signed"

  base_address = "${var.console_fqdn}"
  ca_cert_pem  = "${module.kube_certs.ca_cert_pem}"
  ca_key_alg   = "${module.kube_certs.ca_key_alg}"
  ca_key_pem   = "${module.kube_certs.ca_key_pem}"
}

module "identity_certs" {
  source = "./identity/self-signed"

  ca_cert_pem = "${module.kube_certs.ca_cert_pem}"
  ca_key_alg  = "${module.kube_certs.ca_key_alg}"
  ca_key_pem  = "${module.kube_certs.ca_key_pem}"
}
