module "bootkube" {
  source         = "../../../modules/bootkube"
  cloud_provider = ""

  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"
  oidc_issuer_url    = "https://${var.tectonic_cluster_name}.${var.tectonic_base_domain}:32000/identity"

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

  etcd_servers = ["http://127.0.0.1:2379"]
}

module "tectonic" {
  source   = "../../../modules/tectonic"
  platform = "aws"

  domain             = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}:32000"
  kube_apiserver_url = "https://${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}:443"

  # Platform-independent variables wiring, do not modify.
  container_images = "${var.tectonic_container_images}"
  versions         = "${var.tectonic_versions}"

  license     = "${var.tectonic_license}"
  pull_secret = "${var.tectonic_pull_secret}"

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
  ingress_kind      = "NodePort"
}

module "dns" {
  source = "./dns"

  cluster_name = "${var.tectonic_cluster_name}"
  base_domain  = "${var.tectonic_base_domain}"

  etcd_records = ["${openstack_compute_instance_v2.etcd_node.*.access_ip_v4}"]

  master_records = ["${openstack_compute_instance_v2.master_node.*.access_ip_v4}"]
  master_count   = "${var.tectonic_master_count}"

  worker_records = ["${openstack_compute_instance_v2.worker_node.*.access_ip_v4}"]
  worker_count   = "${var.tectonic_worker_count}"

  tectonic_console_records = ["${openstack_compute_instance_v2.worker_node.*.access_ip_v4}"]
  tectonic_api_records     = ["${openstack_compute_instance_v2.master_node.*.access_ip_v4}"]
}

module "etcd" {
  source = "../../../modules/openstack/etcd"

  count            = "1"
  cluster_name     = "${var.tectonic_cluster_name}"
  core_public_keys = ["${module.secrets.core_public_key_openssh}"]
}

data "null_data_source" "local" {
  inputs = {
    kube_image_url = "${element(split(":", var.tectonic_container_images["hyperkube"]), 0)}"
    kube_image_tag = "${element(split(":", var.tectonic_container_images["hyperkube"]), 1)}"
  }
}

module "nodes" {
  source = "../../../modules/openstack/nodes"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
nameserver 8.8.8.8
nameserver 8.8.4.4
EOF

  kubeconfig_content           = "${module.bootkube.kubeconfig}"
  etcd_fqdns                   = ["${var.tectonic_cluster_name}-etc.${var.tectonic_base_domain}"]
  cluster_name                 = "${var.tectonic_cluster_name}"
  master_count                 = "${var.tectonic_master_count}"
  worker_count                 = "${var.tectonic_master_count}"
  kube_image_url               = "${data.null_data_source.local.outputs.kube_image_url}"
  kube_image_tag               = "${data.null_data_source.local.outputs.kube_image_tag}"
  tectonic_versions            = "${var.tectonic_versions}"
  tectonic_kube_dns_service_ip = "${var.tectonic_kube_dns_service_ip}"

  core_public_keys = ["${module.secrets.core_public_key_openssh}"]
}

module "secrets" {
  source       = "../../../modules/openstack/secrets"
  cluster_name = "${var.tectonic_cluster_name}"
}
