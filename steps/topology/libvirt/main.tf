provider "libvirt" {
  uri = "${var.tectonic_libvirt_uri}"
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

# Set up the cluster domain name
# This is currently limited to the first worker, due to an issue with net-update, even though libvirt supports multiple a-records
resource "null_resource" "console_dns" {
  provisioner "local-exec" {
    command = "virsh -c ${var.tectonic_libvirt_uri} net-update ${libvirt_network.tectonic_net.name} add dns-host \"<host ip='${var.tectonic_libvirt_worker_ips[0]}'><hostname>${var.tectonic_cluster_name}</hostname></host>\" --live --config"
  }
}
