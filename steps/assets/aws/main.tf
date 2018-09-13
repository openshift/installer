module assets_base {
  source = "../base"

  cloud_provider = "aws"
  ingress_kind   = "haproxy-router"

  tectonic_admin_email          = "${var.tectonic_admin_email}"
  tectonic_admin_password       = "${var.tectonic_admin_password}"
  tectonic_admin_ssh_key        = "${var.tectonic_admin_ssh_key}"
  tectonic_base_domain          = "${var.tectonic_base_domain}"
  tectonic_cluster_cidr         = "${var.tectonic_cluster_cidr}"
  tectonic_cluster_id           = "${var.tectonic_cluster_id}"
  tectonic_cluster_name         = "${var.tectonic_cluster_name}"
  tectonic_container_images     = "${var.tectonic_container_images}"
  tectonic_image_re             = "${var.tectonic_image_re}"
  tectonic_kubelet_debug_config = "${var.tectonic_kubelet_debug_config}"
  tectonic_networking           = "${var.tectonic_networking}"
  tectonic_platform             = "${var.tectonic_platform}"
  tectonic_pull_secret          = "${var.tectonic_pull_secret}"
  tectonic_service_cidr         = "${var.tectonic_service_cidr}"
  tectonic_update_channel       = "${var.tectonic_update_channel}"
  tectonic_versions             = "${var.tectonic_versions}"
  aws_worker_ign_config         = "${file("worker.ign")}"
}
