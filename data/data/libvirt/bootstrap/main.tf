resource "libvirt_volume" "bootstrap" {
  name           = "bootstrap"
  base_volume_id = "${var.base_volume_id}"
}

resource "libvirt_ignition" "bootstrap" {
  name    = "bootstrap.ign"
  content = "${var.ignition}"
}

resource "libvirt_domain" "bootstrap" {
  name = "bootstrap"

  memory = "2048"

  vcpu = "2"

  coreos_ignition = "${libvirt_ignition.bootstrap.id}"

  disk {
    volume_id = "${libvirt_volume.bootstrap.id}"
  }

  console {
    type        = "pty"
    target_port = 0
  }

  network_interface {
    network_id = "${var.network_id}"
    hostname   = "${var.cluster_name}-bootstrap"
    addresses  = "${var.addresses}"
  }
}
