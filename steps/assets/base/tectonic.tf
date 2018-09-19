locals {
  ingress_internal_fqdn = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  api_internal_fqdn     = "${var.tectonic_cluster_name}-api.${var.tectonic_base_domain}"
}

data "template_file" "etcd_hostname_list" {
  count    = "${var.tectonic_master_count}"
  template = "${var.tectonic_cluster_name}-etcd-${count.index}.${var.tectonic_base_domain}"
}

module "bootkube" {
  source = "../../../modules/bootkube"

  cluster_name       = "${var.tectonic_cluster_name}"
  kube_apiserver_url = "https://${local.api_internal_fqdn}:6443"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"

  service_cidr = "${var.tectonic_service_cidr}"

  pull_secret = "${var.tectonic_pull_secret}"

  admin_cert_pem                  = "${local.admin_cert_pem}"
  admin_key_pem                   = "${local.admin_key_pem}"
  aggregator_ca_cert_pem          = "${local.aggregator_ca_cert_pem}"
  aggregator_ca_key_pem           = "${local.aggregator_ca_key_pem}"
  apiserver_cert_pem              = "${local.apiserver_cert_pem}"
  apiserver_key_pem               = "${local.apiserver_key_pem}"
  apiserver_proxy_cert_pem        = "${local.apiserver_proxy_cert_pem}"
  apiserver_proxy_key_pem         = "${local.apiserver_proxy_key_pem}"
  etcd_ca_cert_pem                = "${local.etcd_ca_cert_pem}"
  etcd_client_cert_pem            = "${local.etcd_client_cert_pem}"
  etcd_client_key_pem             = "${local.etcd_client_key_pem}"
  kube_ca_cert_pem                = "${local.kube_ca_cert_pem}"
  kube_ca_key_pem                 = "${local.kube_ca_key_pem}"
  kubelet_cert_pem                = "${local.kubelet_cert_pem}"
  kubelet_key_pem                 = "${local.kubelet_key_pem}"
  clusterapi_ca_cert_pem          = "${local.clusterapi_ca_cert_pem}"
  clusterapi_ca_key_pem           = "${local.clusterapi_ca_key_pem}"
  oidc_ca_cert                    = "${local.oidc_ca_cert}"
  openshift_apiserver_cert_pem    = "${local.openshift_apiserver_cert_pem}"
  openshift_apiserver_key_pem     = "${local.openshift_apiserver_key_pem}"
  root_ca_cert_pem                = "${local.root_ca_cert_pem}"
  service_serving_ca_cert_pem     = "${local.service_serving_ca_cert_pem}"
  service_serving_ca_key_pem      = "${local.service_serving_ca_key_pem}"
  mcs_cert_pem                    = "${local.mcs_cert_pem}"
  mcs_key_pem                     = "${local.mcs_key_pem}"
  service_account_public_key_pem  = "${local.service_account_public_key_pem}"
  service_account_private_key_pem = "${local.service_account_private_key_pem}"

  etcd_endpoints = "${data.template_file.etcd_hostname_list.*.rendered}"

  worker_ign_config = "${var.aws_worker_ign_config}"

  libvirt_tls_ca_pem   = "${var.libvirt_tls_ca_pem}"
  libvirt_tls_cert_pem = "${var.libvirt_tls_cert_pem}"
  libvirt_tls_key_pem  = "${var.libvirt_tls_key_pem}"
}

module "tectonic" {
  source = "../../../modules/tectonic"

  base_address = "${local.ingress_internal_fqdn}"

  # Platform-independent variables wiring, do not modify.
  container_images      = "${var.tectonic_container_images}"
  container_base_images = "${var.tectonic_container_base_images}"
  versions              = "${var.tectonic_versions}"

  pull_secret = "${var.tectonic_pull_secret}"

  update_channel = "${var.tectonic_update_channel}"
  update_app_id  = "${var.tectonic_update_app_id}"
  update_server  = "${var.tectonic_update_server}"

  ingress_ca_cert_pem = "${local.ingress_ca_cert_pem}"
  ingress_cert_pem    = "${local.ingress_cert_pem}"
  ingress_key_pem     = "${local.ingress_key_pem}"
  ingress_bundle_pem  = "${join("", list(local.ingress_cert_pem, local.ingress_key_pem, local.ingress_ca_cert_pem))}"

  platform     = "${var.tectonic_platform}"
  ingress_kind = "${var.ingress_kind}"
}
