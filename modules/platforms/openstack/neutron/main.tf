module "dns" {
  source = "./dns"

  cluster_name = "${var.tectonic_cluster_name}"
  base_domain  = "${var.tectonic_base_domain}"

  tectonic_console_records = ["${module.network.worker_floating_ips}"]
  tectonic_api_records     = ["${module.network.master_floating_ips}"]
}

module "etcd" {
  source = "./../etcd"

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

module "master" {
  source = "./../master"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
nameserver 8.8.8.8
nameserver 8.8.4.4
EOF

  kubeconfig_content = "${file("${var.tectonic_assets_dir}/auth/kubeconfig")}"
  etcd_fqdns         = ["${openstack_compute_instance_v2.etcd_node.*.access_ip_v4}"]
  cluster_name       = "${var.tectonic_cluster_name}"
  count              = "${var.tectonic_master_count}"
  kube_image_url     = "${data.null_data_source.local.outputs.kube_image_url}"
  kube_image_tag     = "${data.null_data_source.local.outputs.kube_image_tag}"

  core_public_keys = ["${module.secrets.core_public_key_openssh}"]
}

module "bootkube" {
  source = "./../../../bootkube"

  trigger_ids      = ["${openstack_compute_instance_v2.master_node.*.id}"]
  assets_dir       = "${var.tectonic_assets_dir}"
  core_private_key = "${module.secrets.core_private_key_pem}"
  hosts            = ["${module.network.master_floating_ips}"]
}

module "worker" {
  source = "./../worker"

  resolv_conf_content = <<EOF
search ${var.tectonic_base_domain}
nameserver 8.8.8.8
nameserver 8.8.4.4
EOF

  count              = "${var.tectonic_worker_count}"
  cluster_name       = "${var.tectonic_cluster_name}"
  kubeconfig_content = "${file("${var.tectonic_assets_dir}/auth/kubeconfig")}"
  kube_image_url     = "${data.null_data_source.local.outputs.kube_image_url}"
  kube_image_tag     = "${data.null_data_source.local.outputs.kube_image_tag}"

  etcd_fqdns       = ["${openstack_compute_instance_v2.etcd_node.*.access_ip_v4}"]
  core_public_keys = ["${module.secrets.core_public_key_openssh}"]
}

module "secrets" {
  source       = "./../secrets"
  cluster_name = "${var.tectonic_cluster_name}"
}

module "network" {
  source = "./network"

  master_count        = "${var.tectonic_master_count}"
  worker_count        = "${var.tectonic_worker_count}"
  external_gateway_id = "${var.tectonic_openstack_external_gateway_id}"
  cluster_name        = "${var.tectonic_cluster_name}"
}
