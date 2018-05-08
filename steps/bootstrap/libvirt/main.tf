provider "libvirt" {
  uri = "qemu:///system"
}

locals {
  master_count = 1 # TODO: merge this with the master step
}

resource "libvirt_volume" "master" {
  count = "${local.master_count}"

  name           = "master${count.index}"
  base_volume_id = "${local.libvirt_base_volume_id}"
}

resource "libvirt_ignition" "master" {
  count = "${local.master_count}"

  name    = "master${count.index}.ign"
  content = "${local.ignition_bootstrap}"
}

resource "libvirt_domain" "master" {
  count = "${local.master_count}"

  name = "master${count.index}"

  memory          = "${var.tectonic_libvirt_master_memory}"
  coreos_ignition = "${element(libvirt_ignition.master.*.id,count.index)}"

  disk {
    volume_id = "${element(libvirt_volume.master.*.id, count.index)}"
  }

  network_interface {
    network_id = "${local.libvirt_network_id}"
    hostname   = "${var.tectonic_cluster_name}-master-${count.index}"
    addresses  = ["${var.tectonic_libvirt_master_ips[count.index]}"]
  }
}
