provider "libvirt" {
  uri = "${var.tectonic_libvirt_uri}"
}

resource "libvirt_volume" "master" {
  count = "${var.tectonic_master_count}"

  name           = "master${count.index}"
  base_volume_id = "${local.libvirt_base_volume_id}"
}

resource "libvirt_ignition" "master" {
  name    = "master.ign"
  content = "${file(format("%s/%s", path.cwd, var.tectonic_ignition_master))}"
}

resource "libvirt_domain" "master" {
  count = "${var.tectonic_master_count}"

  name = "master${count.index}"

  memory = "2048"
  vcpu   = "2"

  coreos_ignition = "${libvirt_ignition.master.id}"

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
