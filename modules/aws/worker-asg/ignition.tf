resource "ignition_config" "worker" {
  files = [
    "${ignition_file.kubeconfig.id}",
    "${ignition_file.kubelet-env.id}",
    "${ignition_file.max-user-watches.id}",
    "${ignition_file.etcd-endpoints.id}",
    "${ignition_file.ca-cert.id}",
    "${ignition_file.client-cert.id}",
    "${ignition_file.client-key.id}",
  ]

  systemd = [
    "${ignition_systemd_unit.etcd-member.id}",
    "${ignition_systemd_unit.docker.id}",
    "${ignition_systemd_unit.locksmithd.id}",
    "${ignition_systemd_unit.kubelet-worker.id}",
    "${ignition_systemd_unit.wait-for-dns.id}",
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

data "template_file" "kubelet-worker" {
  template = "${file("${path.module}/resources/worker-kubelet.service")}"

  vars {
    cluster_dns = "${var.tectonic_kube_dns_service_ip}"
  }
}

resource "ignition_systemd_unit" "kubelet-worker" {
  name    = "kubelet.service"
  enable  = true
  content = "${data.template_file.kubelet-worker.rendered}"
}

data "template_file" "etcd-member" {
  template = "${file("${path.module}/resources/etcd-member.service")}"

  vars {
    version   = "${var.tectonic_versions["etcd"]}"
    endpoints = "${join(",",module.etcd.endpoints)}"
  }
}

resource "ignition_systemd_unit" "bootkube" {
  name   = "bootkube.service"
  enable = true

  content = <<EOF
[Unit]
Description=Bootstrap a Kubernetes control plane with a temp api-server
[Service]
Type=simple
WorkingDirectory=/opt/bootkube
ExecStart=/opt/bootkube/assets/bootkube-start
EOF
}

resource "ignition_systemd_unit" "wait-for-dns" {
  name   = "wait-for-dns.service"
  enable = true

  content = <<EOF
[Unit]
Description=Wait for DNS entries
Wants=systemd-resolved.service
Before=kubelet.service
[Service]
Type=oneshot
RemainAfterExit=true
ExecStart=/bin/sh -c 'while ! /usr/bin/grep '^[^#[:space:]]' /etc/resolv.conf \u003e /dev/null; do sleep 1; done'
[Install]
RequiredBy=kubelet.service
EOF
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
    content = "${module.bootkube.kubeconfig}"
  }
}

resource "ignition_file" "kubelet-env" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubelet.env"
  mode       = "420"

  content {
    content = <<EOF
KUBELET_ACI=${var.kube_image_url}
KUBELET_VERSION=${var.kube_image_tag}
EOF
  }
}

resource "ignition_file" "opt-bootkube" {
  filesystem = "root"
  path       = "/opt/bootkube/.empty"
  mode       = "420"

  content {
    content = ""
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

resource "ignition_systemd_unit" "tectonic" {
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
