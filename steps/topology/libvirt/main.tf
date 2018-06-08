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

  coreos_qcow_path = "${var.tectonic_coreos_qcow_path}"
}

locals {
  first_worker_ip = "${cidrhost(var.tectonic_libvirt_ip_range, var.tectonic_libvirt_first_ip_worker)}"
}

# Set up the cluster domain name
# This is currently limited to the first worker, due to an issue with net-update, even though libvirt supports multiple a-records
resource "null_resource" "console_dns" {
  provisioner "local-exec" {
    command = "virsh -c qemu:///system net-update ${var.tectonic_libvirt_network_name} add dns-host \"<host ip='${local.first_worker_ip}'><hostname>${var.tectonic_cluster_name}</hostname></host>\" --live --config"
  }
}
