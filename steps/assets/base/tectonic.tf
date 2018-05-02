locals {
  ingress_internal_fqdn = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  api_internal_fqdn     = "${var.tectonic_cluster_name}-api.${var.tectonic_base_domain}"
}

data "template_file" "etcd_hostname_list" {
  count    = "${var.etcd_count}"
  template = "${var.tectonic_cluster_name}-etcd-${count.index}.${var.tectonic_base_domain}"
}

module "bootkube" {
  source = "../../../modules/bootkube"

  cluster_name       = "${var.tectonic_cluster_name}"
  kube_apiserver_url = "https://${local.api_internal_fqdn}:443"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"

  service_cidr = "${var.tectonic_service_cidr}"

  pull_secret_path = "${pathexpand(var.tectonic_pull_secret_path)}"

  admin_cert_pem           = "${local.admin_cert_pem}"
  admin_key_pem            = "${local.admin_key_pem}"
  aggregator_ca_cert_pem   = "${local.aggregator_ca_cert_pem}"
  apiserver_cert_pem       = "${local.apiserver_cert_pem}"
  apiserver_key_pem        = "${local.apiserver_key_pem}"
  apiserver_proxy_cert_pem = "${local.apiserver_proxy_cert_pem}"
  apiserver_proxy_key_pem  = "${local.apiserver_proxy_key_pem}"
  etcd_ca_cert_pem         = "${local.etcd_ca_cert_pem}"
  etcd_client_cert_pem     = "${local.etcd_client_cert_pem}"
  etcd_client_key_pem      = "${local.etcd_client_key_pem}"
  kube_ca_cert_pem         = "${local.kube_ca_cert_pem}"
  kube_ca_key_pem          = "${local.kube_ca_key_pem}"
  kubelet_cert_pem         = "${local.kubelet_cert_pem}"
  kubelet_key_pem          = "${local.kubelet_key_pem}"
  oidc_ca_cert             = "${local.oidc_ca_cert}"
  root_ca_cert_pem         = "${local.root_ca_cert_pem}"

  etcd_endpoints = "${data.template_file.etcd_hostname_list.*.rendered}"
}

module "tectonic" {
  source = "../../../modules/tectonic"

  base_address = "${local.ingress_internal_fqdn}"

  # Platform-independent variables wiring, do not modify.
  container_images      = "${var.tectonic_container_images}"
  container_base_images = "${var.tectonic_container_base_images}"
  versions              = "${var.tectonic_versions}"

  license_path     = "${pathexpand(var.tectonic_license_path)}"
  pull_secret_path = "${pathexpand(var.tectonic_pull_secret_path)}"

  admin_email = "${var.tectonic_admin_email}"

  update_channel = "${var.tectonic_update_channel}"
  update_app_id  = "${var.tectonic_update_app_id}"
  update_server  = "${var.tectonic_update_server}"

  ca_generated = "${var.tectonic_ca_cert == "" ? false : true}"

  ingress_ca_cert_pem = "${local.ingress_ca_cert_pem}"
  ingress_cert_pem    = "${local.ingress_cert_pem}"
  ingress_key_pem     = "${local.ingress_key_pem}"

  identity_client_ca_cert  = "${local.identity_client_ca_cert}"
  identity_client_cert_pem = "${local.identity_client_cert_pem}"
  identity_client_key_pem  = "${local.identity_client_key_pem}"
  identity_server_ca_cert  = "${local.identity_server_ca_cert}"
  identity_server_cert_pem = "${local.identity_server_cert_pem}"
  identity_server_key_pem  = "${local.identity_server_key_pem}"

  platform     = "${var.tectonic_platform}"
  ingress_kind = "${var.ingress_kind}"
}
