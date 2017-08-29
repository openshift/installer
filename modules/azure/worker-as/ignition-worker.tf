data "ignition_config" "worker" {
  files = [
    "${data.ignition_file.kubeconfig.id}",
    "${var.ign_kubelet_env_id}",
    "${var.ign_azure_udev_rules_id}",
    "${var.ign_max_user_watches_id}",
    "${data.ignition_file.cloud-provider-config.id}",
  ]

  systemd = [
    "${var.ign_docker_dropin_id}",
    "${var.ign_locksmithd_service_id}",
    "${var.ign_kubelet_service_id}",
    "${var.ign_tx_off_service_id}",
  ]

  users = [
    "${data.ignition_user.core.id}",
  ]
}

data "ignition_file" "kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"
  mode       = 0644

  content {
    content = "${var.kubeconfig_content}"
  }
}

data "ignition_file" "cloud-provider-config" {
  filesystem = "root"
  path       = "/etc/kubernetes/cloud/config"
  mode       = 0600

  content {
    content = "${var.cloud_provider_config}"
  }
}

data "ignition_user" "core" {
  name = "core"

  ssh_authorized_keys = [
    "${file(var.public_ssh_key)}",
  ]
}
