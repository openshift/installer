module "bootkube" {
  source         = "../../../modules/bootkube"
  cloud_provider = ""

  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"
  oidc_issuer_url    = "https://${var.tectonic_cluster_name}.${var.tectonic_base_domain}/identity"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"
  versions         = "${var.tectonic_versions}"

  ca_cert    = "${var.tectonic_ca_cert}"
  ca_key     = "${var.tectonic_ca_key}"
  ca_key_alg = "${var.tectonic_ca_key_alg}"

  service_cidr = "${var.tectonic_service_cidr}"
  cluster_cidr = "${var.tectonic_cluster_cidr}"

  advertise_address = "0.0.0.0"
  anonymous_auth    = "false"

  oidc_username_claim = "email"
  oidc_groups_claim   = "groups"
  oidc_client_id      = "tectonic-kubectl"

  etcd_endpoints   = ["${aws_route53_record.etc_a_nodes.*.fqdn}"]
  etcd_ca_cert     = "${var.tectonic_etcd_ca_cert_path}"
  etcd_client_cert = "${var.tectonic_etcd_client_cert_path}"
  etcd_client_key  = "${var.tectonic_etcd_client_key_path}"
  etcd_tls_enabled = "${var.tectonic_etcd_tls_enabled}"

  etcd_cert_dns_names = [
    "${var.tectonic_cluster_name}-etcd-0.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-1.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-2.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-3.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-4.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-5.${var.tectonic_base_domain}",
    "${var.tectonic_cluster_name}-etcd-6.${var.tectonic_base_domain}",
  ]

  experimental_enabled = "${var.tectonic_experimental}"

  master_count = "${var.tectonic_master_count}"
}

module "tectonic" {
  source   = "../../../modules/tectonic"
  platform = "aws"

  base_address       = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"
  versions         = "${var.tectonic_versions}"

  license_path     = "${var.tectonic_vanilla_k8s ? "/dev/null" : pathexpand(var.tectonic_license_path)}"
  pull_secret_path = "${var.tectonic_vanilla_k8s ? "/dev/null" : pathexpand(var.tectonic_pull_secret_path)}"

  admin_email         = "${var.tectonic_admin_email}"
  admin_password_hash = "${var.tectonic_admin_password_hash}"

  update_channel = "${var.tectonic_update_channel}"
  update_app_id  = "${var.tectonic_update_app_id}"
  update_server  = "${var.tectonic_update_server}"

  ca_generated = "${module.bootkube.ca_cert == "" ? false : true}"
  ca_cert      = "${module.bootkube.ca_cert}"
  ca_key_alg   = "${module.bootkube.ca_key_alg}"
  ca_key       = "${module.bootkube.ca_key}"

  console_client_id = "tectonic-console"
  kubectl_client_id = "tectonic-kubectl"
  ingress_kind      = "HostPort"
  experimental      = "${var.tectonic_experimental}"
  master_count      = "${var.tectonic_master_count}"
  stats_url         = "${var.tectonic_stats_url}"
}

module "etcd" {
  source = "../../../modules/openstack/etcd"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
nameserver 8.8.8.8
nameserver 8.8.4.4
EOF

  base_domain           = "${var.tectonic_base_domain}"
  cluster_name          = "${var.tectonic_cluster_name}"
  container_image       = "${var.tectonic_container_images["etcd"]}"
  core_public_keys      = ["${module.secrets.core_public_key_openssh}"]
  tectonic_experimental = "${var.tectonic_experimental}"
  tls_enabled           = "${var.tectonic_etcd_tls_enabled}"

  tls_ca_crt_pem     = "${module.bootkube.etcd_ca_crt_pem}"
  tls_client_crt_pem = "${module.bootkube.etcd_client_crt_pem}"
  tls_client_key_pem = "${module.bootkube.etcd_client_key_pem}"
  tls_peer_crt_pem   = "${module.bootkube.etcd_peer_crt_pem}"
  tls_peer_key_pem   = "${module.bootkube.etcd_peer_key_pem}"

  instance_count = "${var.tectonic_etcd_count}"
}

data "null_data_source" "local" {
  inputs = {
    kube_image_url = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$1")}"
    kube_image_tag = "${replace(var.tectonic_container_images["hyperkube"],var.tectonic_image_re,"$2")}"
  }
}

module "master_nodes" {
  source = "../../../modules/openstack/nodes"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
nameserver 8.8.8.8
nameserver 8.8.4.4
EOF

  kubeconfig_content           = "${module.bootkube.kubeconfig}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  instance_count               = "${var.tectonic_master_count}"
  kube_image_url               = "${data.null_data_source.local.outputs.kube_image_url}"
  kube_image_tag               = "${data.null_data_source.local.outputs.kube_image_tag}"
  tectonic_kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  core_public_keys             = ["${module.secrets.core_public_key_openssh}"]
  bootkube_service             = "${module.bootkube.systemd_service}"
  tectonic_service             = "${module.tectonic.systemd_service}"
  hostname_infix               = "master"
  node_labels                  = "node-role.kubernetes.io/master"
  node_taints                  = "node-role.kubernetes.io/master=:NoSchedule"
  tectonic_experimental        = "${var.tectonic_experimental}"
}

module "worker_nodes" {
  source = "../../../modules/openstack/nodes"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
nameserver 8.8.8.8
nameserver 8.8.4.4
EOF

  kubeconfig_content           = "${module.bootkube.kubeconfig}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  instance_count               = "${var.tectonic_worker_count}"
  kube_image_url               = "${data.null_data_source.local.outputs.kube_image_url}"
  kube_image_tag               = "${data.null_data_source.local.outputs.kube_image_tag}"
  tectonic_kube_dns_service_ip = "${module.bootkube.kube_dns_service_ip}"
  core_public_keys             = ["${module.secrets.core_public_key_openssh}"]
  bootkube_service             = ""
  tectonic_service             = ""
  hostname_infix               = "worker"
  node_labels                  = "node-role.kubernetes.io/node"
  node_taints                  = ""
}

module "secrets" {
  source       = "../../../modules/openstack/secrets"
  cluster_name = "${var.tectonic_cluster_name}"
}

module "secgroups" {
  source                = "../../../modules/openstack/secgroups"
  cluster_name          = "${var.tectonic_cluster_name}"
  tectonic_experimental = "${var.tectonic_experimental}"
}
