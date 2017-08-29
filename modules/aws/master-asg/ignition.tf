data "ignition_config" "main" {
  files = [
    "${var.ign_max_user_watches_id}",
    "${var.ign_s3_puller_id}",
    "${data.ignition_file.init_assets.id}",
    "${data.ignition_file.detect_master.id}",
  ]

  systemd = [
    "${var.ign_docker_dropin_id}",
    "${var.ign_locksmithd_service_id}",
    "${var.ign_kubelet_service_id}",
    "${var.ign_s3_kubelet_env_service_id}",
    "${data.ignition_systemd_unit.init_assets.id}",
    "${data.ignition_systemd_unit.bootkube.id}",
    "${data.ignition_systemd_unit.tectonic.id}",
  ]
}

data "ignition_file" "detect_master" {
  filesystem = "root"
  path       = "/opt/detect-master.sh"
  mode       = 0755

  content {
    content = "${file("${path.module}/resources/detect-master.sh")}"
  }
}

data "template_file" "init_assets" {
  template = "${file("${path.module}/resources/init-assets.sh")}"

  vars {
    cluster_name       = "${var.cluster_name}"
    awscli_image       = "${var.container_images["awscli"]}"
    assets_s3_location = "${var.assets_s3_location}"
    kubelet_image_url  = "${replace(var.container_images["hyperkube"],var.image_re,"$1")}"
    kubelet_image_tag  = "${replace(var.container_images["hyperkube"],var.image_re,"$2")}"
  }
}

data "ignition_file" "init_assets" {
  filesystem = "root"
  path       = "/opt/tectonic/init-assets.sh"
  mode       = 0755

  content {
    content = "${data.template_file.init_assets.rendered}"
  }
}

data "ignition_systemd_unit" "init_assets" {
  name    = "init-assets.service"
  enable  = "${var.assets_s3_location != "" ? true : false}"
  content = "${file("${path.module}/resources/services/init-assets.service")}"
}

data "ignition_systemd_unit" "bootkube" {
  name    = "bootkube.service"
  content = "${var.bootkube_service}"
}

data "ignition_systemd_unit" "tectonic" {
  name    = "tectonic.service"
  enable  = "${var.tectonic_service_disabled == 0 ? true : false}"
  content = "${var.tectonic_service}"
}
