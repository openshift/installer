data "ignition_config" "worker" {
  files = [
    "${data.ignition_file.kubeconfig.id}",
    "${data.ignition_file.kubelet-env.id}",
    "${data.ignition_file.max-user-watches.id}",
  ]

  systemd = [
    "${data.ignition_systemd_unit.etcd-member.id}",
    "${data.ignition_systemd_unit.docker.id}",
    "${data.ignition_systemd_unit.locksmithd.id}",
    "${data.ignition_systemd_unit.kubelet-worker.id}",
  ]

  users = [
    "${data.ignition_user.core.id}",
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

  dropin = [
    {
      name    = "40-etcd-lock.conf"
      content = "[Service]\nEnvironment=REBOOT_STRATEGY=etcd-lock\n"
    },
  ]
}

data "template_file" "kubelet-worker" {
  template = "${file("${path.module}/resources/worker-kubelet.service")}"

  vars {
    node_label     = "${var.kubelet_node_label}"
    cloud_provider = "${var.cloud_provider}"
    cluster_dns    = "${var.tectonic_kube_dns_service_ip}"
  }
}

data "ignition_systemd_unit" "kubelet-worker" {
  name    = "kubelet.service"
  enable  = true
  content = "${data.template_file.kubelet-worker.rendered}"
}

data "template_file" "etcd-member" {
  template = "${file("${path.module}/resources/etcd-member.service")}"

  vars {
    version   = "${var.tectonic_versions["etcd"]}"
    endpoints = "${join(",",formatlist("%s:2379",var.etcd_endpoints))}"
  }
}

data "ignition_systemd_unit" "etcd-member" {
  name   = "etcd-member.service"
  enable = true

  dropin = [
    {
      name    = "40-etcd-gateway.conf"
      content = "${data.template_file.etcd-member.rendered}"
    },
  ]
}

data "ignition_file" "kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"
  mode       = "420"

  content {
    content = "${var.kubeconfig_content}"
  }
}

data "ignition_file" "kubelet-env" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubelet.env"
  mode       = "420"

  content {
    content = <<EOF
KUBELET_IMAGE_URL="${var.kube_image_url}"
KUBELET_IMAGE_TAG="${var.kube_image_tag}"
EOF
  }
}

data "ignition_file" "max-user-watches" {
  filesystem = "root"
  path       = "/etc/sysctl.d/max-user-watches.conf"
  mode       = "420"

  content {
    content = "fs.inotify.max_user_watches=16184"
  }
}

data "ignition_systemd_unit" "tectonic" {
  name   = "tectonic.service"
  enable = true

  content = <<EOF
[Unit]
Description=Bootstrap a Tectonic cluster
[Service]
Type=oneshot
WorkingDirectory=/opt/tectonic
ExecStart=/usr/bin/bash /opt/tectonic/bootkube.sh
ExecStart=/usr/bin/bash /opt/tectonic/tectonic.sh kubeconfig tectonic
EOF
}

data "ignition_user" "core" {
  name = "core"

  ssh_authorized_keys = [
    "${file(var.public_ssh_key)}",
  ]
}
