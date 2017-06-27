data "ignition_config" "node" {
  count = "${var.instance_count}"

  users = [
    "${data.ignition_user.core.id}",
  ]

  files = [
    "${data.ignition_file.max-user-watches.id}",
    "${data.ignition_file.node_hostname.*.id[count.index]}",
    "${data.ignition_file.kubelet-env.id}",
  ]

  systemd = [
    "${data.ignition_systemd_unit.docker.id}",
    "${data.ignition_systemd_unit.locksmithd.id}",
    "${data.ignition_systemd_unit.kubelet.id}",
    "${data.ignition_systemd_unit.kubelet-env.id}",
    "${data.ignition_systemd_unit.bootkube.id}",
    "${data.ignition_systemd_unit.tectonic.id}",
    "${data.ignition_systemd_unit.vmtoolsd_member.id}",
  ]

  networkd = [
    "${data.ignition_networkd_unit.vmnetwork.*.id[count.index]}",
  ]
}

data "ignition_user" "core" {
  name                = "core"
  ssh_authorized_keys = ["${var.core_public_keys}"]
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
    cluster_dns_ip    = "${var.kube_dns_service_ip}"
    node_label        = "${var.kubelet_node_label}"
    node_taints_param = "${var.kubelet_node_taints != "" ? "--register-with-taints=${var.kubelet_node_taints}" : ""}"
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
    kube_version_image_url = "${replace(var.container_images["kube_version"],var.image_re,"$1")}"
    kube_version_image_tag = "${replace(var.container_images["kube_version"],var.image_re,"$2")}"
    kubelet_image_url      = "${replace(var.container_images["hyperkube"],var.image_re,"$1")}"
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
  mode       = 0644

  content {
    content = "fs.inotify.max_user_watches=16184"
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

data "ignition_systemd_unit" "vmtoolsd_member" {
  name   = "vmtoolsd.service"
  enable = true

  content = <<EOF
  [Unit]
  Description=VMware Tools Agent
  Documentation=http://open-vm-tools.sourceforge.net/
  ConditionVirtualization=vmware
  [Service]
  ExecStartPre=/usr/bin/ln -sfT /usr/share/oem/vmware-tools /etc/vmware-tools
  ExecStart=/usr/share/oem/bin/vmtoolsd
  TimeoutStopSec=5
EOF
}

data "ignition_networkd_unit" "vmnetwork" {
  count = "${var.instance_count}"
  name  = "00-ens192.network"

  content = <<EOF
  [Match]
  Name=ens192
  [Network]
  DNS=${var.dns_server}
  Address=${var.ip_address["${count.index}"]}
  Gateway=${var.gateway}
  UseDomains=yes
  Domains=${var.base_domain}
EOF
}

data "ignition_file" "node_hostname" {
  count      = "${var.instance_count}"
  path       = "/etc/hostname"
  mode       = 0644
  filesystem = "root"

  content {
    content = "${var.hostname["${count.index}"]}"
  }
}

data "ignition_file" "kubelet-env" {
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
