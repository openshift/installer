provider "aws" {
  region  = "${var.tectonic_aws_region}"
  profile = "${var.tectonic_aws_profile}"
  version = "1.8.0"

  assume_role {
    role_arn     = "${var.tectonic_aws_installer_role == "" ? "" : "${var.tectonic_aws_installer_role}"}"
    session_name = "TECTONIC_INSTALLER_${var.tectonic_cluster_name}"
  }
}

locals {
  ingress_internal_fqdn = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  api_internal_fqdn     = "${var.tectonic_cluster_name}-api.${var.tectonic_base_domain}"
  bucket_name           = "${var.tectonic_cluster_name}-tnc.${var.tectonic_base_domain}"
  bucket_assets_key     = "assets.zip"
}

data "aws_availability_zones" "azs" {}

data "template_file" "etcd_hostname_list" {
  count    = "${var.tectonic_etcd_count > 0 ? var.tectonic_etcd_count : length(data.aws_availability_zones.azs.names) == 5 ? 5 : 3}"
  template = "${var.tectonic_cluster_name}-etcd-${count.index}.${var.tectonic_base_domain}"
}

module "bootkube" {
  source = "../../modules/bootkube"

  cluster_name       = "${var.tectonic_cluster_name}"
  kube_apiserver_url = "https://${local.api_internal_fqdn}:443"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"

  service_cidr = "${var.tectonic_service_cidr}"

  pull_secret_path = "${pathexpand(var.tectonic_pull_secret_path)}"

  oidc_ca_cert             = "${module.ingress_certs.ca_cert_pem}"
  root_ca_cert_pem         = "${module.ca_certs.root_ca_cert_pem}"
  aggregator_ca_cert_pem   = "${module.ca_certs.aggregator_ca_cert_pem}"
  apiserver_cert_pem       = "${module.kube_certs.apiserver_cert_pem}"
  apiserver_key_pem        = "${module.kube_certs.apiserver_key_pem}"
  apiserver_proxy_cert_pem = "${module.kube_certs.apiserver_proxy_cert_pem}"
  apiserver_proxy_key_pem  = "${module.kube_certs.apiserver_proxy_key_pem}"
  etcd_ca_cert_pem         = "${module.ca_certs.etcd_ca_cert_pem}"
  etcd_ca_key_pem          = "${module.ca_certs.etcd_ca_key_pem}"
  etcd_client_cert_pem     = "${module.etcd_certs.etcd_client_cert_pem}"
  etcd_client_key_pem      = "${module.etcd_certs.etcd_client_key_pem}"

  kube_ca_cert_pem = "${module.ca_certs.kube_ca_cert_pem}"
  kube_ca_key_pem  = "${module.ca_certs.kube_ca_key_pem}"
  admin_cert_pem   = "${module.kube_certs.admin_cert_pem}"
  admin_key_pem    = "${module.kube_certs.admin_key_pem}"

  etcd_endpoints = "${data.template_file.etcd_hostname_list.*.rendered}"
}

module "tectonic" {
  source   = "../../modules/tectonic"
  platform = "aws"

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

  ingress_ca_cert_pem = "${module.ingress_certs.ca_cert_pem}"
  ingress_cert_pem    = "${module.ingress_certs.cert_pem}"
  ingress_key_pem     = "${module.ingress_certs.key_pem}"

  identity_client_ca_cert  = "${module.ca_certs.root_ca_cert_pem}"
  identity_client_cert_pem = "${module.identity_certs.client_cert_pem}"
  identity_client_key_pem  = "${module.identity_certs.client_key_pem}"
  identity_server_ca_cert  = "${module.ca_certs.kube_ca_cert_pem}"
  identity_server_cert_pem = "${module.identity_certs.server_cert_pem}"
  identity_server_key_pem  = "${module.identity_certs.server_key_pem}"

  ingress_kind = "NodePort"
}
