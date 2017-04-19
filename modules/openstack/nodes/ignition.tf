resource "ignition_config" "node" {
  count = "${var.instance_count}"

  users = [
    "${ignition_user.core.id}",
  ]

  files = [
    "${ignition_file.kubeconfig.id}",
    "${ignition_file.kubelet-env.id}",
    "${ignition_file.max-user-watches.id}",
    "${ignition_file.resolv_conf.id}",
    "${ignition_file.hostname.*.id[count.index]}",
  ]

  systemd = [
    "${ignition_systemd_unit.etcd-member.id}",
    "${ignition_systemd_unit.docker.id}",
    "${ignition_systemd_unit.locksmithd.id}",
    "${ignition_systemd_unit.kubelet.id}",
    "${ignition_systemd_unit.bootkube.id}",
    "${ignition_systemd_unit.tectonic.id}",
  ]
}

resource "ignition_file" "resolv_conf" {
  path       = "/etc/resolv.conf"
  mode       = 0644
  uid        = 0
  filesystem = "root"

  content {
    content = "${var.resolv_conf_content}"
  }
}

resource "ignition_user" "core" {
  name                = "core"
  ssh_authorized_keys = ["${var.core_public_keys}"]
}

resource "ignition_file" "hostname" {
  count      = "${var.instance_count}"
  path       = "/etc/hostname"
  mode       = 0644
  uid        = 0
  filesystem = "root"

  content {
    content = "${var.cluster_name}-${var.hostname_infix}-${count.index}"
  }
}

resource "ignition_systemd_unit" "docker" {
  name   = "docker.service"
  enable = true
}

resource "ignition_systemd_unit" "locksmithd" {
  name = "locksmithd.service"

  dropin = [
    {
      name    = "40-etcd-lock.conf"
      content = "[Service]\nEnvironment=REBOOT_STRATEGY=etcd-lock\n"
    },
  ]
}

data "template_file" "kubelet" {
  template = "${file("${path.module}/resources/kubelet.service")}"

  vars {
    cluster_dns       = "${var.tectonic_kube_dns_service_ip}"
    node_labels       = "${var.node_labels}"
    node_taints_param = "${var.node_taints != "" ? "--register-with-taints=${var.node_taints}" : ""}"
  }
}

resource "ignition_systemd_unit" "kubelet" {
  name    = "kubelet.service"
  enable  = true
  content = "${data.template_file.kubelet.rendered}"
}

data "template_file" "etcd-member" {
  template = "${file("${path.module}/resources/etcd-member.service")}"

  vars {
    version   = "${var.tectonic_versions["etcd"]}"
    endpoints = "${join(",", formatlist("%s:2379", var.etcd_fqdns))}"
  }
}

resource "ignition_systemd_unit" "etcd-member" {
  name   = "etcd-member.service"
  enable = true

  dropin = [
    {
      name    = "40-etcd-gateway.conf"
      content = "${data.template_file.etcd-member.rendered}"
    },
  ]
}

resource "ignition_file" "kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"
  mode       = "420"

  content {
    content = "${var.kubeconfig_content}"
  }
}

resource "ignition_file" "kubelet-env" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubelet.env"
  mode       = "420"

  content {
    content = <<EOF
KUBELET_IMAGE_URL=${var.kube_image_url}
KUBELET_IMAGE_TAG="${var.kube_image_tag}"
EOF
  }
}

resource "ignition_file" "max-user-watches" {
  filesystem = "root"
  path       = "/etc/sysctl.d/max-user-watches.conf"
  mode       = "420"

  content {
    content = "fs.inotify.max_user_watches=16184"
  }
}

resource "ignition_systemd_unit" "bootkube" {
  name    = "bootkube.service"
  content = "${var.bootkube_service}"
}

resource "ignition_systemd_unit" "tectonic" {
  name    = "tectonic.service"
  content = "${var.tectonic_service}"
}
