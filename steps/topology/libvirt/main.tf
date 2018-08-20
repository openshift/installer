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

  dns_host = ["${flatten(list(
    data.libvirt_network_dns_host_template.masters.*.rendered,
    data.libvirt_network_dns_host_template.etcds.*.rendered,
    data.libvirt_network_dns_host_template.workers.*.rendered,
  ))}"]
}

module "libvirt_base_volume" {
  source = "../../../modules/libvirt/volume"

  coreos_qcow_path = "${var.tectonic_coreos_qcow_path}"
}

locals {
  "hostnames" = [
    "${var.tectonic_cluster_name}-api",
    "${var.tectonic_cluster_name}-tnc",
  ]
}

data "libvirt_network_dns_host_template" "masters" {
  count = "${var.tectonic_master_count * length(local.hostnames)}"

  ip = "${var.tectonic_libvirt_master_ips[count.index / length(local.hostnames)]}"

  hostname = "${local.hostnames[count.index % length(local.hostnames)]}"
}

data "libvirt_network_dns_host_template" "etcds" {
  count = "${var.tectonic_etcd_count}"

  ip = "${var.tectonic_libvirt_etcd_ips[count.index]}"

  hostname = "${var.tectonic_cluster_name}-etcd-${count.index}"
}

data "libvirt_network_dns_host_template" "workers" {
  count = "${var.tectonic_worker_count}"

  ip = "${var.tectonic_libvirt_worker_ips[count.index]}"

  hostname = "${var.tectonic_cluster_name}"
}
