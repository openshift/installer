data "ignition_config" "main" {
  files = [
    "${data.ignition_file.max-user-watches.id}",
    "${data.ignition_file.s3-puller.id}",
    "${data.ignition_file.init-assets.id}",
    "${data.ignition_file.detect-master.id}",
  ]

  systemd = [
    "${data.ignition_systemd_unit.docker.id}",
    "${data.ignition_systemd_unit.locksmithd.id}",
    "${data.ignition_systemd_unit.kubelet.id}",
    "${data.ignition_systemd_unit.kubelet-env.id}",
    "${data.ignition_systemd_unit.init-assets.id}",
    "${data.ignition_systemd_unit.bootkube.id}",
    "${data.ignition_systemd_unit.tectonic.id}",
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

data "template_file" "kubelet" {
  template = "${file("${path.module}/resources/services/kubelet.service")}"

  vars {
    cluster_dns_ip         = "${var.kube_dns_service_ip}"
    node_label             = "${var.kubelet_node_label}"
    node_taints_param      = "${var.kubelet_node_taints != "" ? "--register-with-taints=${var.kubelet_node_taints}" : ""}"
    kubeconfig_s3_location = "${var.kubeconfig_s3_location}"
  }
}

data "ignition_systemd_unit" "kubelet" {
  name    = "kubelet.service"
  enable  = true
  content = "${data.template_file.kubelet.rendered}"
}

data "template_file" "kubelet-env" {
  template = "${file("${path.module}/resources/services/kubelet-env.service")}"

  vars {
    kube_version_image_url = "${element(split(":", var.container_images["kube_version"]), 0)}"
    kube_version_image_tag = "${element(split(":", var.container_images["kube_version"]), 1)}"
    kubelet_image_url      = "${element(split(":", var.container_images["hyperkube"]), 0)}"
    kubeconfig_s3_location = "${var.kubeconfig_s3_location}"
  }
}

data "ignition_systemd_unit" "kubelet-env" {
  name    = "kubelet-env.service"
  enable  = true
  content = "${data.template_file.kubelet-env.rendered}"
}

data "ignition_file" "max-user-watches" {
  filesystem = "root"
  path       = "/etc/sysctl.d/max-user-watches.conf"
  mode       = "420"

  content {
    content = "fs.inotify.max_user_watches=16184"
  }
}

data "ignition_file" "s3-puller" {
  filesystem = "root"
  path       = "/opt/s3-puller.sh"
  mode       = "555"

  content {
    content = "${file("${path.module}/resources/s3-puller.sh")}"
  }
}

data "ignition_file" "detect-master" {
  filesystem = "root"
  path       = "/opt/detect-master.sh"
  mode       = "555"

  content {
    content = "${file("${path.module}/resources/detect-master.sh")}"
  }
}

data "template_file" "init-assets" {
  template = "${file("${path.module}/resources/init-assets.sh")}"

  vars {
    cluster_name       = "${var.cluster_name}"
    awscli_image       = "${var.container_images["awscli"]}"
    assets_s3_location = "${var.assets_s3_location}"
    kubelet_image_url  = "${element(split(":", var.container_images["hyperkube"]), 0)}"
    kubelet_image_tag  = "${element(split(":", var.container_images["hyperkube"]), 1)}"
  }
}

data "ignition_file" "init-assets" {
  filesystem = "root"
  path       = "/opt/tectonic/init-assets.sh"
  mode       = "555"

  content {
    content = "${data.template_file.init-assets.rendered}"
  }
}

data "ignition_systemd_unit" "init-assets" {
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
