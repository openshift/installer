provider "libvirt" {
  uri = "${var.tectonic_libvirt_uri}"
}

locals {
  master_count = "${var.tectonic_bootstrap == "true" ? 1 : var.tectonic_master_count}"
}

resource "libvirt_volume" "master" {
  count = "${local.master_count}"

  name           = "master${count.index}"
  base_volume_id = "${local.libvirt_base_volume_id}"
}

# The first master node should be booted with the bootstrap ignition configuration
resource "libvirt_ignition" "master_bootstrap" {
  name    = "master-bootstrap.ign"
  content = "${local.ignition_bootstrap}"
}

# Ignition for the remaining masters
resource "libvirt_ignition" "master" {
  name    = "master.ign"
  content = "${file("${path.cwd}/${var.tectonic_ignition_master}")}"
}

resource "libvirt_domain" "master" {
  count = "${local.master_count}"

  name = "master${count.index}"

  memory = "${var.tectonic_libvirt_master_memory}"

  # Override ignition for the first (bootstrap) node. It can't be re-ignited,
  # but that's okay for us
  coreos_ignition = "${count.index == 0 ? libvirt_ignition.master_bootstrap.id : libvirt_ignition.master.id}"

  disk {
    volume_id = "${element(libvirt_volume.master.*.id, count.index)}"
  }

  console {
    type        = "pty"
    target_port = 0
  }

  network_interface {
    network_id = "${local.libvirt_network_id}"
    hostname   = "${var.tectonic_cluster_name}-master-${count.index}"
    addresses  = ["${var.tectonic_libvirt_master_ips[count.index]}"]
  }
}
