data "ignition_config" "main" {
  files = [
    "${var.ign_max_user_watches_id}",
    "${var.ign_gcs_puller_id}",
    "${data.ignition_file.init_assets.id}",
    "${var.ign_installer_kubelet_env_id}",
  ]

  systemd = ["${compact(list(
    var.ign_docker_dropin_id,
    var.ign_locksmithd_service_id,
    var.ign_kubelet_service_id,
    var.ign_k8s_node_bootstrap_service_id,
    var.ign_init_assets_service_id,
    var.ign_bootkube_service_id,
    var.ign_tectonic_service_id,
    var.ign_bootkube_path_unit_id,
    var.ign_tectonic_path_unit_id
   ))}"]
}

data "template_file" "init_assets" {
  template = "${file("${path.module}/resources/init-assets.sh")}"

  vars {
    cluster_name        = "${var.cluster_name}"
    assets_gcs_location = "${var.assets_gcs_location}"
    kubelet_image_url   = "${replace(var.container_images["hyperkube"],var.image_re,"$1")}"
    kubelet_image_tag   = "${replace(var.container_images["hyperkube"],var.image_re,"$2")}"
  }
}

data "ignition_file" "init_assets" {
  filesystem = "root"
  path       = "/opt/init-assets.sh"
  mode       = 0755

  content {
    content = "${data.template_file.init_assets.rendered}"
  }
}
