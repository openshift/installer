provider "libvirt" {
  uri = "${var.tectonic_libvirt_uri}"
}

module "volume" {
  source = "./volume"

  cluster_name = "${var.tectonic_cluster_name}"
  image        = "${var.tectonic_os_image}"
}

module "bootstrap" {
  source = "./bootstrap"

  addresses      = ["${var.tectonic_libvirt_bootstrap_ip}"]
  base_volume_id = "${module.volume.coreos_base_volume_id}"
  cluster_name   = "${var.tectonic_cluster_name}"
  ignition       = "${var.ignition_bootstrap}"
  network_id     = "${libvirt_network.tectonic_net.id}"
}

resource "libvirt_volume" "master" {
  count          = "${var.tectonic_master_count}"
  name           = "${var.tectonic_cluster_name}-master-${count.index}"
  base_volume_id = "${module.volume.coreos_base_volume_id}"
}

resource "libvirt_ignition" "master" {
  name    = "${var.tectonic_cluster_name}-master.ign"
  content = "${var.ignition_master}"
}

resource "libvirt_network" "tectonic_net" {
  name = "${var.tectonic_cluster_name}"

  mode   = "nat"
  bridge = "${var.tectonic_libvirt_network_if}"

  domain = "${var.tectonic_base_domain}"

  addresses = [
    "${var.tectonic_libvirt_ip_range}",
  ]

  dns = [{
    local_only = true

    srvs = ["${flatten(list(
      data.libvirt_network_dns_srv_template.etcd_cluster.*.rendered,
    ))}"]

    hosts = ["${flatten(list(
      data.libvirt_network_dns_host_template.bootstrap.*.rendered,
      data.libvirt_network_dns_host_template.masters.*.rendered,
      data.libvirt_network_dns_host_template.etcds.*.rendered,
    ))}"]
  }]

  autostart = true
}

resource "libvirt_domain" "master" {
  count = "${var.tectonic_master_count}"

  name = "${var.tectonic_cluster_name}-master-${count.index}"

  memory = "3072"
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

data "libvirt_network_dns_host_template" "bootstrap" {
  count    = "${var.bootstrap_dns ? 1 : 0}"
  ip       = "${var.tectonic_libvirt_bootstrap_ip}"
  hostname = "${var.tectonic_cluster_name}-api"
}

data "libvirt_network_dns_host_template" "masters" {
  count    = "${var.tectonic_master_count}"
  ip       = "${var.tectonic_libvirt_master_ips[count.index]}"
  hostname = "${var.tectonic_cluster_name}-api"
}

data "libvirt_network_dns_host_template" "etcds" {
  count    = "${var.tectonic_master_count}"
  ip       = "${var.tectonic_libvirt_master_ips[count.index]}"
  hostname = "${var.tectonic_cluster_name}-etcd-${count.index}"
}

data "libvirt_network_dns_srv_template" "etcd_cluster" {
  count    = "${var.tectonic_master_count}"
  service  = "etcd-server-ssl"
  protocol = "tcp"
  domain   = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  port     = 2380
  weight   = 10
  target   = "${var.tectonic_cluster_name}-etcd-${count.index}.${var.tectonic_base_domain}"
}
