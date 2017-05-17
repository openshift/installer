module "bootkube" {
  source         = "../../../modules/bootkube"
  cloud_provider = ""

  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"
  oidc_issuer_url    = "https://${var.tectonic_cluster_name}.${var.tectonic_base_domain}/identity"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"

  ca_cert    = "${var.tectonic_ca_cert}"
  ca_key     = "${var.tectonic_ca_key}"
  ca_key_alg = "${var.tectonic_ca_key_alg}"

  service_cidr = "${var.tectonic_service_cidr}"
  cluster_cidr = "${var.tectonic_cluster_cidr}"

  kube_apiserver_service_ip = "${var.tectonic_kube_apiserver_service_ip}"
  kube_dns_service_ip       = "${var.tectonic_kube_dns_service_ip}"

  advertise_address = "0.0.0.0"
  anonymous_auth    = "false"

  oidc_username_claim = "email"
  oidc_groups_claim   = "groups"
  oidc_client_id      = "tectonic-kubectl"

  etcd_endpoints   = ["${openstack_compute_instance_v2.etcd_node.*.access_ip_v4}"]
  etcd_ca_cert     = "${var.tectonic_etcd_ca_cert_path}"
  etcd_client_cert = "${var.tectonic_etcd_client_cert_path}"
  etcd_client_key  = "${var.tectonic_etcd_client_key_path}"
  etcd_service_ip  = "${var.tectonic_kube_etcd_service_ip}"
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

  master_count = "${var.tectonic_master_count}"
}

data "null_data_source" "local" {
  inputs = {
    kube_image_url = "${element(split(":", var.tectonic_container_images["hyperkube"]), 0)}"
    kube_image_tag = "${element(split(":", var.tectonic_container_images["hyperkube"]), 1)}"
  }
}

module "etcd" {
  source = "../../../modules/openstack/etcd"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
nameserver 8.8.8.8
nameserver 8.8.4.4
EOF

  base_domain      = "${var.tectonic_base_domain}"
  cluster_name     = "${var.tectonic_cluster_name}"
  container_image  = "${var.tectonic_container_images["etcd"]}"
  core_public_keys = ["${module.secrets.core_public_key_openssh}"]
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
  tectonic_kube_dns_service_ip = "${var.tectonic_kube_dns_service_ip}"
  core_public_keys             = ["${module.secrets.core_public_key_openssh}"]
  bootkube_service             = "${module.bootkube.systemd_service}"
  tectonic_service             = "${module.tectonic.systemd_service}"
  hostname_infix               = "master"
  node_labels                  = "node-role.kubernetes.io/master"
  node_taints                  = "node-role.kubernetes.io/master=:NoSchedule"
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
  tectonic_kube_dns_service_ip = "${var.tectonic_kube_dns_service_ip}"
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
