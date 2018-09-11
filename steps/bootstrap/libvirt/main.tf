provider "libvirt" {
  uri = "${var.tectonic_libvirt_uri}"
}

module "bootstrap" {
  source = "../../../modules/libvirt/bootstrap"

  addresses      = ["${var.tectonic_libvirt_bootstrap_ip}"]
  base_volume_id = "${local.libvirt_base_volume_id}"
  cluster_name   = "${var.tectonic_cluster_name}"
  ignition       = "${local.ignition_bootstrap}"
  network_id     = "${local.libvirt_network_id}"
}
