provider "libvirt" {
  uri = "qemu:///system" #XXX fixme
}

module "defaults" {
  source = "../../../modules/libvirt/target-defaults"

  etcd_count = "${var.tectonic_etcd_count}"
}

resource "libvirt_volume" "etcd" {
  count          = "${module.defaults.etcd_count}"
  name           = "etcd${count.index}"
  base_volume_id = "${local.libvirt_base_volume_id}"
}

resource "libvirt_ignition" "etcd" {
  count   = "${module.defaults.etcd_count}"
  name    = "etcd${count.index}.ign"
  content = "${local.ignition[count.index]}"
}

resource "libvirt_domain" "etcd" {
  count = "${module.defaults.etcd_count}"

  name            = "etcd${count.index}"
  memory          = "${var.tectonic_libvirt_etcd_memory}"
  coreos_ignition = "${element(libvirt_ignition.etcd.*.id,count.index)}"

  disk {
    volume_id = "${element(libvirt_volume.etcd.*.id, count.index)}"
  }

  network_interface {
    network_id = "${local.libvirt_network_id}"
    hostname   = "${var.tectonic_cluster_name}-etcd-${count.index}"
    addresses  = ["${cidrhost(var.tectonic_libvirt_ip_range, var.tectonic_libvirt_first_ip_etcd + count.index)}"]
  }
}
