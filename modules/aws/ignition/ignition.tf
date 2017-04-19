resource "ignition_config" "main" {
  files = [
    "${ignition_file.max-user-watches.id}",
    "${ignition_file.s3-puller.id}",
    "${ignition_file.init-assets.id}",
  ]

  systemd = [
    "${ignition_systemd_unit.etcd-member.id}",
    "${ignition_systemd_unit.docker.id}",
    "${ignition_systemd_unit.locksmithd.id}",
    "${ignition_systemd_unit.kubelet.id}",
    "${ignition_systemd_unit.init-assets.id}",
    "${ignition_systemd_unit.bootkube.id}",
    "${ignition_systemd_unit.tectonic.id}",
  ]
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
  template = "${file("${path.module}/resources/services/kubelet.service")}"

  vars {
    aci                    = "${element(split(":", var.container_images["hyperkube"]), 0)}"
    version                = "${element(split(":", var.container_images["hyperkube"]), 1)}"
    cluster_dns_ip         = "${var.kube_dns_service_ip}"
    node_label             = "${var.kubelet_node_label}"
    node_taints_param      = "${var.kubelet_node_taints != "" ? "--register-with-taints=${var.kubelet_node_taints}" : ""}"
    kubeconfig_s3_location = "${var.kubeconfig_s3_location}"
  }
}

resource "ignition_systemd_unit" "kubelet" {
  name    = "kubelet.service"
  enable  = true
  content = "${data.template_file.kubelet.rendered}"
}

data "template_file" "etcd-member" {
  template = "${file("${path.module}/resources/services/etcd-member.service")}"

  vars {
    image     = "${var.container_images["etcd"]}"
    endpoints = "${join(",", formatlist("%s:2379", var.etcd_endpoints))}"
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

resource "ignition_file" "max-user-watches" {
  filesystem = "root"
  path       = "/etc/sysctl.d/max-user-watches.conf"
  mode       = "420"

  content {
    content = "fs.inotify.max_user_watches=16184"
  }
}

resource "ignition_file" "s3-puller" {
  filesystem = "root"
  path       = "/opt/s3-puller.sh"
  mode       = "555"

  content {
    content = "${file("${path.module}/resources/s3-puller.sh")}"
  }
}

data "template_file" "init-assets" {
  template = "${file("${path.module}/resources/init-assets.sh")}"

  vars {
    awscli_image       = "${var.container_images["awscli"]}"
    assets_s3_location = "${var.assets_s3_location}"
  }
}

resource "ignition_file" "init-assets" {
  filesystem = "root"
  path       = "/opt/tectonic/init-assets.sh"
  mode       = "555"

  content {
    content = "${data.template_file.init-assets.rendered}"
  }
}

resource "ignition_systemd_unit" "init-assets" {
  name    = "init-assets.service"
  enable  = "${var.assets_s3_location != "" ? true : false}"
  content = "${file("${path.module}/resources/services/init-assets.service")}"
}

resource "ignition_systemd_unit" "bootkube" {
  name    = "bootkube.service"
  content = "${var.bootkube_service}"
}

resource "ignition_systemd_unit" "tectonic" {
  name    = "tectonic.service"
  enable  = "${var.tectonic_service_disabled == 0 ? true : false}"
  content = "${var.tectonic_service}"
}
