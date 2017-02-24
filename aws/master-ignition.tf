resource "ignition_config" "master" {
  files = [
    "${ignition_file.etcd-endpoints.id}",
    "${ignition_file.kubeconfig.id}",
  ]

  systemd = [
    "${ignition_systemd_unit.docker-master.id}",
    "${ignition_systemd_unit.locksmithd-master.id}",
    "${ignition_systemd_unit.kubelet-master.id}",
  ]
}

resource "ignition_systemd_unit" "docker-master" {
  name   = "docker.service"
  enable = true
}

resource "ignition_systemd_unit" "locksmithd-master" {
  name = "locksmithd.service"

  dropin = [
    {
      name    = "40-etcd-lock.conf"
      content = "[Service]\nEnvironment=REBOOT_STRATEGY=etcd-lock\nEnvironment=LOCKSMITHD_ENDPOINT=${join(",",module.etcd.endpoints)}"
    },
  ]
}

resource "ignition_systemd_unit" "kubelet-master" {
  name    = "kubelet.service"
  enable  = true
  content = "${file("${path.module}/resources/master-kubelet.service")}"
}

resource "ignition_file" "etcd-endpoints" {
  filesystem = "root"
  path       = "/etc/kubernetes/etcd-endpoints.env"

  content {
    content = "TECTONIC_ETCD_ENDPOINTS=${join(",",module.etcd.endpoints)}"
  }
}

resource "ignition_file" "kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"

  content {
    content = "${file("${path.module}/../assets/auth/kubeconfig")}"
  }
}
