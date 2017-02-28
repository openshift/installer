resource "ignition_config" "master" {
  files = [
    "${ignition_file.etcd-endpoints.id}",
    "${ignition_file.kubeconfig.id}",
    "${ignition_file.kubelet-env.id}",
    "${ignition_file.opt-bootkube.id}",
    "${ignition_file.ca-cert.id}",
    "${ignition_file.client-cert.id}",
    "${ignition_file.client-key.id}",
  ]

  systemd = [
    "${ignition_systemd_unit.etcd-member.id}",
    "${ignition_systemd_unit.docker.id}",
    "${ignition_systemd_unit.locksmithd.id}",
    "${ignition_systemd_unit.kubelet-master.id}",
    "${ignition_systemd_unit.wait-for-dns.id}",
    "${ignition_systemd_unit.bootkube.id}",
  ]
}

resource "ignition_config" "worker" {
  files = [
    "${ignition_file.etcd-endpoints.id}",
    "${ignition_file.kubeconfig.id}",
    "${ignition_file.kubelet-env.id}",
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

resource "ignition_systemd_unit" "kubelet-master" {
  name    = "kubelet.service"
  enable  = true
  content = "${file("${path.module}/resources/master-kubelet.service")}"
}

resource "ignition_systemd_unit" "kubelet-worker" {
  name    = "kubelet.service"
  enable  = true
  content = "${file("${path.module}/resources/worker-kubelet.service")}"
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
      name = "40-etcd-gateway.conf"

      content = <<EOF
[Service]
Environment="ETCD_IMAGE_TAG=v3.1.0"
EnvironmentFile=/etc/kubernetes/etcd-endpoints.env
ExecStart=
ExecStart=/usr/lib/coreos/etcd-wrapper gateway start \
    --listen-addr=127.0.0.1:2379 \
    --endpoints=$${TECTONIC_ETCD_ENDPOINTS}
EOF
    },
  ]
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
    content = <<EOF
apiVersion: v1
kind: Config
clusters:
- name: local
  cluster:
    server: https://${var.cluster_name}-k8s.${var.tectonic_domain}:443
    certificate-authority: /etc/kubernetes/ssl/ca.pem
users:
- name: kubelet
  user:
    client-certificate: /etc/kubernetes/ssl/client.pem
    client-key: /etc/kubernetes/ssl/client-key.pem
contexts:
- context:
    cluster: local
    user: kubelet
EOF
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

resource "ignition_file" "ca-cert" {
  filesystem = "root"
  path       = "/etc/kubernetes/ssl/ca.pem"
  mode       = "420"

  content {
    content = "${file("${path.root}/../assets/tls/ca.crt")}"
  }
}
