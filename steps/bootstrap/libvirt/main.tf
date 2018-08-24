provider "libvirt" {
  uri = "${var.tectonic_libvirt_uri}"
}

resource "libvirt_volume" "bootstrap" {
  name           = "bootstrap"
  base_volume_id = "${local.libvirt_base_volume_id}"
}

resource "libvirt_ignition" "bootstrap" {
  name    = "bootstrap.ign"
  content = "${local.ignition_bootstrap}"
}

resource "libvirt_domain" "bootstrap" {
  name = "bootstrap"

  memory = "2048"

  vcpu = "2"

  coreos_ignition = "${libvirt_ignition.bootstrap.id}"

  disk {
    volume_id = "${libvirt_volume.bootstrap.id}"
  }

  network_interface {
    network_id = "${local.libvirt_network_id}"
    hostname   = "${var.tectonic_cluster_name}-bootstrap"
    addresses  = ["${var.tectonic_libvirt_bootstrap_ip}"]
  }
}
