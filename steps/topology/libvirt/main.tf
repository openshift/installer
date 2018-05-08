provider "libvirt" {
  uri = "qemu:///system" #XXX fixme
}

# Create the bridge for libvirt
resource "libvirt_network" "tectonic_net" {
  name = "${var.tectonic_libvirt_network_name}"

  mode   = "nat"
  bridge = "${var.tectonic_libvirt_network_if}"

  domain = "${var.tectonic_base_domain}"

  addresses = [
    "${var.tectonic_libvirt_ip_range}",
  ]

  dns_forwarder {
    address = "${var.tectonic_libvirt_resolver}"
  }
}

module "libvirt_base_volume" {
  source = "../../../modules/libvirt/volume"

  coreos_qow_path = "${var.tectonic_coreos_qow_path}"
}
