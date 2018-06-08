# Terraform doesn't support "inheritance"
# So we have to pass all variables down
module assets_base {
  source = "../base"

  cloud_provider = ""
  etcd_count     = "${var.tectonic_etcd_count > 0 ? var.tectonic_etcd_count : 1}"

  ingress_kind = "haproxy-router"

  tectonic_base_domain             = "${var.tectonic_base_domain}"
  tectonic_bootstrap_upgrade_cl    = "${var.tectonic_bootstrap_upgrade_cl}"
  tectonic_cluster_name            = "${var.tectonic_cluster_name}"
  tectonic_container_images        = "${var.tectonic_container_images}"
  tectonic_image_re                = "${var.tectonic_image_re}"
  tectonic_kubelet_debug_config    = "${var.tectonic_kubelet_debug_config}"
  tectonic_cluster_cidr            = "${var.tectonic_cluster_cidr}"
  tectonic_service_cidr            = "${var.tectonic_service_cidr}"
  tectonic_networking              = "${var.tectonic_networking}"
  tectonic_license_path            = "${var.tectonic_license_path}"
  tectonic_pull_secret_path        = "${var.tectonic_pull_secret_path}"
  tectonic_admin_email             = "${var.tectonic_admin_email}"
  tectonic_update_channel          = "${var.tectonic_update_channel}"
  tectonic_platform                = "${var.tectonic_platform}"
  tectonic_versions                = "${var.tectonic_versions}"
  tectonic_admin_password          = "${var.tectonic_admin_password}"
  tectonic_cluster_id              = "${var.tectonic_cluster_id}"
  tectonic_container_linux_channel = "${var.tectonic_container_linux_channel}"
  tectonic_container_linux_version = "${var.tectonic_container_linux_version}"
}

# Removing assets is platform-specific
# But it must be installed in /opt/tectonic/rm-assets.sh
data "ignition_file" "rm_assets_sh" {
  filesystem = "root"
  path       = "/opt/tectonic/rm-assets.sh"
  mode       = "0700"

  content {
    content = "${file("${path.module}/resources/rm-assets.sh")}"
  }
}

data "ignition_systemd_unit" "remap_port" {
  name    = "remap_port.service"
  enabled = true
  content = "${file("${path.module}/resources/remap_port.service")}"
}

data "ignition_user" "core" {
  name = "core"

  ssh_authorized_keys = [
    "${var.tectonic_libvirt_ssh_key}",
  ]
}

data "ignition_config" "bootstrap" {
  files = ["${flatten(list(
    list(
      data.ignition_file.rm_assets_sh.id,
    ),
    module.assets_base.ignition_bootstrap_files,
  ))}"]

  systemd = [
    "${module.assets_base.ignition_bootstrap_systemd}",
    "${data.ignition_systemd_unit.remap_port.id}",
  ]

  users = [
    "${data.ignition_user.core.id}",
  ]
}

data "ignition_config" "etcd" {
  count = "${var.tectonic_etcd_count}"

  users = [
    "${data.ignition_user.core.id}",
  ]

  files = [
    "${module.assets_base.ignition_etcd_files}",
  ]

  append {
    source = "${format("http://${var.tectonic_cluster_name}-tnc.${var.tectonic_base_domain}:49500/config/etcd?etcd_index=%d", count.index)}"
  }
}
