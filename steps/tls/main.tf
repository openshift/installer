locals {
  api_internal_fqdn     = "${var.tectonic_cluster_name}-api.${var.tectonic_base_domain}"
  ingress_internal_fqdn = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
}

module "ca_certs" {
  source = "../../modules/tls/ca"

  root_ca_cert_pem_path = "${var.tectonic_ca_cert}"
  root_ca_key_alg       = "${var.tectonic_ca_key_alg}"
  root_ca_key_pem_path  = "${var.tectonic_ca_key}"
}

module "kube_certs" {
  source = "../../modules/tls/kube"

  kube_ca_cert_pem            = "${module.ca_certs.kube_ca_cert_pem}"
  kube_ca_key_alg             = "${module.ca_certs.kube_ca_key_alg}"
  kube_ca_key_pem             = "${module.ca_certs.kube_ca_key_pem}"
  aggregator_ca_cert_pem      = "${module.ca_certs.aggregator_ca_cert_pem}"
  aggregator_ca_key_alg       = "${module.ca_certs.aggregator_ca_key_alg}"
  aggregator_ca_key_pem       = "${module.ca_certs.aggregator_ca_key_pem}"
  service_serving_ca_cert_pem = "${module.ca_certs.service_serving_ca_cert_pem}"
  service_serving_ca_key_alg  = "${module.ca_certs.service_serving_ca_key_alg}"
  service_serving_ca_key_pem  = "${module.ca_certs.service_serving_ca_key_pem}"
  kube_apiserver_url          = "https://${local.api_internal_fqdn}:443"
  service_cidr                = "${var.tectonic_service_cidr}"
}

module "etcd_certs" {
  source = "../../modules/tls/etcd"

  etcd_ca_cert_pem = "${module.ca_certs.etcd_ca_cert_pem}"
  etcd_ca_key_alg  = "${module.ca_certs.etcd_ca_key_alg}"
  etcd_ca_key_pem  = "${module.ca_certs.etcd_ca_key_pem}"
}

module "ingress_certs" {
  source = "../../modules/tls/ingress"

  base_address = "${local.ingress_internal_fqdn}"
  ca_cert_pem  = "${module.ca_certs.kube_ca_cert_pem}"
  ca_key_alg   = "${module.ca_certs.kube_ca_key_alg}"
  ca_key_pem   = "${module.ca_certs.kube_ca_key_pem}"
}
