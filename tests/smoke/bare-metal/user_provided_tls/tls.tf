module "ca_certs" {
  source = "../../modules/tls/ca/user-provided"

  root_ca_cert_pem_path       = "../../tests/smoke/bare-metal/user_provided_tls/certs/ca/root_ca.crt"
  root_ca_key_pem_path        = "../../tests/smoke/bare-metal/user_provided_tls/certs/ca/root_ca.key"
  etcd_ca_cert_pem_path       = "../../tests/smoke/bare-metal/user_provided_tls/certs/ca/etcd_ca.crt"
  etcd_ca_key_pem_path        = "../../tests/smoke/bare-metal/user_provided_tls/certs/ca/etcd_ca.key"
  kube_ca_cert_pem_path       = "../../tests/smoke/bare-metal/user_provided_tls/certs/ca/kube_ca.crt"
  kube_ca_key_pem_path        = "../../tests/smoke/bare-metal/user_provided_tls/certs/ca/kube_ca.key"
  aggregator_ca_cert_pem_path = "../../tests/smoke/bare-metal/user_provided_tls/certs/ca/aggregator_ca.crt"
  aggregator_ca_key_pem_path  = "../../tests/smoke/bare-metal/user_provided_tls/certs/ca/aggregator_ca.key"
}

module "kube_certs" {
  source = "../../modules/tls/kube"

  kube_ca_cert_pem       = "${module.ca_certs.kube_ca_cert_pem}"
  kube_ca_key_alg        = "${module.ca_certs.kube_ca_key_alg}"
  kube_ca_key_pem        = "${module.ca_certs.kube_ca_key_pem}"
  aggregator_ca_cert_pem = "${module.ca_certs.aggregator_ca_cert_pem}"
  aggregator_ca_key_alg  = "${module.ca_certs.aggregator_ca_key_alg}"
  aggregator_ca_key_pem  = "${module.ca_certs.aggregator_ca_key_pem}"
  kube_apiserver_url     = "https://${var.tectonic_metal_controller_domain}:6443"
  service_cidr           = "${var.tectonic_service_cidr}"
}

module "etcd_certs" {
  source = "../../modules/tls/etcd"

  etcd_ca_cert_pem    = "${module.ca_certs.etcd_ca_cert_pem}"
  etcd_ca_key_alg     = "${module.ca_certs.etcd_ca_key_alg}"
  etcd_ca_key_pem     = "${module.ca_certs.etcd_ca_key_pem}"
  service_cidr        = "${var.tectonic_service_cidr}"
  etcd_cert_dns_names = "${var.tectonic_metal_controller_domains}"
}

module "ingress_certs" {
  source = "../../modules/tls/ingress/user-provided"

  ca_cert_pem_path = "../../tests/smoke/bare-metal/user_provided_tls/certs/ingress/ca.crt"
  cert_pem_path    = "../../tests/smoke/bare-metal/user_provided_tls/certs/ingress/ingress.crt"
  key_pem_path     = "../../tests/smoke/bare-metal/user_provided_tls/certs/ingress/ingress.key"
}

module "identity_certs" {
  source = "../../modules/tls/identity"

  kube_ca_cert_pem = "${module.ca_certs.kube_ca_cert_pem}"
  kube_ca_key_alg  = "${module.ca_certs.kube_ca_key_alg}"
  kube_ca_key_pem  = "${module.ca_certs.kube_ca_key_pem}"
}
