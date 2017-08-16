data "ignition_config" "master" {
  files = [
    "${data.ignition_file.kubeconfig.id}",
    "${data.ignition_file.kubelet_env.id}",
    "${module.azure_udev-rules.udev-rules_id}",
    "${data.ignition_file.max_user_watches.id}",
    "${data.ignition_file.cloud_provider_config.id}",
  ]

  systemd = [
    "${data.ignition_systemd_unit.docker.id}",
    "${data.ignition_systemd_unit.locksmithd.id}",
    "${data.ignition_systemd_unit.kubelet_master.id}",
    "${data.ignition_systemd_unit.tectonic.id}",
    "${data.ignition_systemd_unit.bootkube.id}",
    "${module.net_ignition.tx-off_id}",
  ]

  users = [
    "${data.ignition_user.core.id}",
  ]
}

data "ignition_user" "core" {
  name = "core"

  ssh_authorized_keys = [
    "${file(var.public_ssh_key)}",
  ]
}

data "ignition_systemd_unit" "docker" {
  name   = "docker.service"
  enable = true

  dropin = [
    {
      name    = "10-dockeropts.conf"
      content = "[Service]\nEnvironment=\"DOCKER_OPTS=--log-opt max-size=50m --log-opt max-file=3\"\n"
    },
  ]
}

data "ignition_systemd_unit" "locksmithd" {
  name = "locksmithd.service"
  mask = true
}

data "template_file" "kubelet-master" {
  template = "${file("${path.module}/resources/master-kubelet.service")}"

  vars {
    node_label        = "${var.kubelet_node_label}"
    node_taints_param = "${var.kubelet_node_taints != "" ? "--register-with-taints=${var.kubelet_node_taints}" : ""}"
    cni_bin_dir_flag  = "${var.kubelet_cni_bin_dir != "" ? "--cni-bin-dir=${var.kubelet_cni_bin_dir}" : ""}"
    cloud_provider    = "${var.cloud_provider}"
    cluster_dns       = "${var.tectonic_kube_dns_service_ip}"
  }
}

data "ignition_systemd_unit" "kubelet_master" {
  name    = "kubelet.service"
  enable  = true
  content = "${data.template_file.kubelet-master.rendered}"
}

data "ignition_file" "kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"
  mode       = 0644

  content {
    content = "${var.kubeconfig_content}"
  }
}

data "ignition_file" "kubelet_env" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubelet.env"
  mode       = 0644

  content {
    content = <<EOF
KUBELET_IMAGE_URL="${var.kube_image_url}"
KUBELET_IMAGE_TAG="${var.kube_image_tag}"
EOF
  }
}

data "ignition_file" "max_user_watches" {
  filesystem = "root"
  path       = "/etc/sysctl.d/max-user-watches.conf"
  mode       = 0644

  content {
    content = "fs.inotify.max_user_watches=16184"
  }
}

data "ignition_file" "cloud_provider_config" {
  filesystem = "root"
  path       = "/etc/kubernetes/cloud/config"
  mode       = 0600

  content {
    content = "${var.cloud_provider_config}"
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

module "net_ignition" {
  source = "../../net/ignition"
}

module "azure_udev-rules" {
  source = "../udev-rules"
}
