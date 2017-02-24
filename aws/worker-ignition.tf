resource "ignition_config" "worker" {
  files = [
    "${ignition_file.etcd-endpoints.id}",
    "${ignition_file.kubeconfig.id}",
  ]

  systemd = [
    "${ignition_systemd_unit.docker.id}",
    "${ignition_systemd_unit.locksmithd.id}",
    "${ignition_systemd_unit.kubelet.id}",
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

resource "ignition_systemd_unit" "kubelet" {
  name    = "kubelet.service"
  enable  = true
  content = "${file("${path.module}/resources/worker-kubelet.service")}"
}
