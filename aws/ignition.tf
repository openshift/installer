resource "ignition_config" "master" {
  files = [
    "${ignition_file.etcd-endpoints.id}",
    "${ignition_file.kubeconfig.id}",
    "${ignition_file.kubelet-env.id}",
  ]

  systemd = [
    "${ignition_systemd_unit.docker.id}",
    "${ignition_systemd_unit.locksmithd.id}",
    "${ignition_systemd_unit.kubelet.id}",
    "${ignition_systemd_unit.wait-for-dns.id}",
    "${ignition_systemd_unit.bootkube.id}",
  ]
}

resource "ignition_config" "worker" {
  files = [
    "${ignition_file.etcd-endpoints.id}",
    "${ignition_file.kubeconfig.id}",
    "${ignition_file.kubelet-env.id}",
  ]

  systemd = [
    "${ignition_systemd_unit.docker.id}",
    "${ignition_systemd_unit.locksmithd.id}",
    "${ignition_systemd_unit.kubelet.id}",
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
      content = "[Service]\nEnvironment=REBOOT_STRATEGY=etcd-lock\nEnvironment=LOCKSMITHD_ENDPOINT=${join(",",module.etcd.endpoints)}"
    },
  ]
}

resource "ignition_systemd_unit" "kubelet" {
  name    = "kubelet.service"
  enable  = true
  content = "${file("${path.module}/resources/master-kubelet.service")}"
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

resource "ignition_file" "etcd-endpoints" {
  filesystem = "root"
  path       = "/etc/kubernetes/etcd-endpoints.env"
  mode       = "420"

  content {
    content = "TECTONIC_ETCD_ENDPOINTS=${join(",",module.etcd.endpoints)}"
  }
}

resource "ignition_file" "kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"
  mode       = "420"

  content {
    content = "${file("${path.module}/../assets/auth/kubeconfig")}"
  }
}

resource "ignition_file" "kubelet-env" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubelet.env"
  mode       = "420"

  content {
    content = <<EOF
KUBELET_ACI=quay.io/coreos/hyperkube
KUBELET_VERSION="${var.kube_version}"
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

resource "ignition_file" "client-ca" {
  filesystem = "root"
  path       = "/etc/kubernetes/ssl/ca.pem"
  mode       = "420"

  content {
    content = "${file("${path.root}/../assets/tls/ca.crt")}"
  }
}

resource "ignition_file" "client-key" {
  filesystem = "root"
  path       = "/etc/kubernetes/ssl/client-key.pem"
  mode       = "420"

  content {
    content = "${file("${path.root}/../assets/tls/kubelet.key")}"
  }
}

resource "ignition_file" "client-cert" {
  filesystem = "root"
  path       = "/etc/kubernetes/ssl/client.pem"
  mode       = "420"

  content {
    content = "${file("${path.root}/../assets/tls/kubelet.crt")}"
  }
}
