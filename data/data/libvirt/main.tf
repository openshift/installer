provider "libvirt" {
  uri = "${var.tectonic_libvirt_uri}"
}

module "libvirt_base_volume" {
  source = "./volume"

  image = "${var.tectonic_os_image}"
}

module "bootstrap" {
  source = "./bootstrap"

  addresses      = ["${var.tectonic_libvirt_bootstrap_ip}"]
  base_volume_id = "${module.libvirt_base_volume.coreos_base_volume_id}"
  cluster_name   = "${var.tectonic_cluster_name}"
  ignition       = "${var.ignition_bootstrap}"
  network_id     = "${libvirt_network.tectonic_net.id}"
}

resource "libvirt_volume" "master" {
  count          = "${var.tectonic_master_count}"
  name           = "master${count.index}"
  base_volume_id = "${module.libvirt_base_volume.coreos_base_volume_id}"
}

resource "libvirt_ignition" "master" {
  count   = "${var.tectonic_master_count}"
  name    = "master-${count.index}.ign"
  content = "${var.ignition_masters[count.index]}"
}

resource "libvirt_ignition" "worker" {
  name    = "worker.ign"
  content = "${var.ignition_worker}"
}

resource "libvirt_network" "tectonic_net" {
  name = "${var.tectonic_libvirt_network_name}"

  mode   = "nat"
  bridge = "${var.tectonic_libvirt_network_if}"

  domain = "${var.tectonic_base_domain}"

  addresses = [
    "${var.tectonic_libvirt_ip_range}",
  ]

  dns = [{
    local_only = true

    hosts = ["${flatten(list(
      data.libvirt_network_dns_host_template.bootstrap.*.rendered,
      data.libvirt_network_dns_host_template.masters.*.rendered,
      data.libvirt_network_dns_host_template.etcds.*.rendered,
      data.libvirt_network_dns_host_template.workers.*.rendered,
    ))}"]
  }]

  autostart = true
}

resource "libvirt_domain" "master" {
  count = "${var.tectonic_master_count}"

  name = "master${count.index}"

  memory = "2048"
  vcpu   = "2"

  coreos_ignition = "${libvirt_ignition.master.*.id[count.index]}"

  disk {
    volume_id = "${element(libvirt_volume.master.*.id, count.index)}"
  }

  console {
    type        = "pty"
    target_port = 0
  }

  network_interface {
    network_id = "${libvirt_network.tectonic_net.id}"
    hostname   = "${var.tectonic_cluster_name}-master-${count.index}"
    addresses  = ["${var.tectonic_libvirt_master_ips[count.index]}"]
  }
}

locals {
  "hostnames" = [
    "${var.tectonic_cluster_name}-api",
  ]
}

data "libvirt_network_dns_host_template" "bootstrap" {
  count    = "${length(local.hostnames)}"
  ip       = "${var.tectonic_libvirt_bootstrap_ip}"
  hostname = "${local.hostnames[count.index]}"
}

data "libvirt_network_dns_host_template" "masters" {
  count    = "${var.tectonic_master_count * length(local.hostnames)}"
  ip       = "${var.tectonic_libvirt_master_ips[count.index / length(local.hostnames)]}"
  hostname = "${local.hostnames[count.index % length(local.hostnames)]}"
}

data "libvirt_network_dns_host_template" "etcds" {
  count    = "${var.tectonic_master_count}"
  ip       = "${var.tectonic_libvirt_master_ips[count.index]}"
  hostname = "${var.tectonic_cluster_name}-etcd-${count.index}"
}

data "libvirt_network_dns_host_template" "workers" {
  count    = "${var.tectonic_worker_count}"
  ip       = "${var.tectonic_libvirt_worker_ips[count.index]}"
  hostname = "${var.tectonic_cluster_name}"
}
