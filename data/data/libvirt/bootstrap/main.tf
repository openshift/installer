resource "libvirt_volume" "bootstrap" {
  name           = "${var.cluster_id}-bootstrap"
  base_volume_id = var.base_volume_id
  pool           = var.pool
}

resource "libvirt_ignition" "bootstrap" {
  name    = "${var.cluster_id}-bootstrap.ign"
  content = var.ignition
  pool    = var.pool
}

resource "libvirt_domain" "bootstrap" {
  name = "${var.cluster_id}-bootstrap"

  memory = var.bootstrap_memory

  vcpu = "2"

  coreos_ignition = libvirt_ignition.bootstrap.id

  disk {
    volume_id = libvirt_volume.bootstrap.id
  }

  console {
    type        = "pty"
    target_port = 0
  }

  cpu = {
    mode = "host-passthrough"
  }

  network_interface {
    network_id = var.network_id
    hostname   = "${var.cluster_id}-bootstrap.${var.cluster_domain}"
    addresses  = var.addresses
  }
}

