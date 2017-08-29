data "ignition_config" "node" {
  count = "${var.instance_count}"

  users = [
    "${data.ignition_user.core.id}",
  ]

  files = [
    "${data.ignition_file.kubeconfig.id}",
    "${var.ign_kubelet_env_id}",
    "${var.ign_max_user_watches_id}",
    "${data.ignition_file.resolv_conf.id}",
    "${data.ignition_file.hostname.*.id[count.index]}",
  ]

  systemd = [
    "${var.ign_docker_dropin_id}",
    "${var.ign_locksmithd_service_id}",
    "${var.ign_kubelet_service_id}",
    "${data.ignition_systemd_unit.bootkube.id}",
    "${data.ignition_systemd_unit.tectonic.id}",
  ]
}

data "ignition_file" "resolv_conf" {
  path       = "/etc/resolv.conf"
  mode       = 0644
  uid        = 0
  filesystem = "root"

  content {
    content = "${var.resolv_conf_content}"
  }
}

data "ignition_user" "core" {
  name                = "core"
  ssh_authorized_keys = ["${var.core_public_keys}"]
}

data "ignition_file" "hostname" {
  count      = "${var.instance_count}"
  path       = "/etc/hostname"
  mode       = 0644
  uid        = 0
  filesystem = "root"

  content {
    content = "${var.cluster_name}-${var.hostname_infix}-${count.index}"
  }
}

data "ignition_file" "kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"
  mode       = 0644

  content {
    content = "${var.kubeconfig_content}"
  }
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
