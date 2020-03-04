provider "libvirt" {
  uri = var.libvirt_uri
}

resource "libvirt_pool" "storage_pool" {
  name = var.cluster_id
  type = "dir"
  path = "${var.libvirt_base_pooldir}/openshift-images/${var.cluster_id}"
}

module "volume" {
  source = "./volume"

  cluster_id = var.cluster_id
  image      = var.os_image
  pool       = libvirt_pool.storage_pool.name
}

module "bootstrap" {
  source = "./bootstrap"

  cluster_domain = var.cluster_domain
  addresses      = [var.libvirt_bootstrap_ip]
  base_volume_id = module.volume.coreos_base_volume_id
  cluster_id     = var.cluster_id
  ignition       = var.ignition_bootstrap
  network_id     = libvirt_network.net.id
  pool           = libvirt_pool.storage_pool.name
}

resource "libvirt_volume" "master" {
  count          = var.master_count
  name           = "${var.cluster_id}-master-${count.index}"
  base_volume_id = module.volume.coreos_base_volume_id
  pool           = libvirt_pool.storage_pool.name
}

resource "libvirt_ignition" "master" {
  name    = "${var.cluster_id}-master.ign"
  content = var.ignition_master
  pool    = libvirt_pool.storage_pool.name
}

resource "libvirt_network" "net" {
  name = var.cluster_id

  mode   = "nat"
  bridge = var.libvirt_network_if

  domain = var.cluster_domain

  addresses = [var.machine_cidr]

  dns {
    local_only = true

    dynamic "srvs" {
      for_each = data.libvirt_network_dns_srv_template.etcd_cluster.*.rendered
      content {
        domain   = srvs.value.domain
        port     = srvs.value.port
        protocol = srvs.value.protocol
        service  = srvs.value.service
        target   = srvs.value.target
        weight   = srvs.value.weight
      }
    }

    dynamic "hosts" {
      for_each = concat(
        data.libvirt_network_dns_host_template.bootstrap.*.rendered,
        data.libvirt_network_dns_host_template.bootstrap_int.*.rendered,
        data.libvirt_network_dns_host_template.masters.*.rendered,
        data.libvirt_network_dns_host_template.masters_int.*.rendered,
        data.libvirt_network_dns_host_template.etcds.*.rendered,
      )
      content {
        hostname = hosts.value.hostname
        ip       = hosts.value.ip
      }
    }
  }

  autostart = true
}

resource "libvirt_domain" "master" {
  count = var.master_count

  name = "${var.cluster_id}-master-${count.index}"

  memory = var.libvirt_master_memory
  vcpu   = var.libvirt_master_vcpu

  coreos_ignition = libvirt_ignition.master.id

  disk {
    volume_id = element(libvirt_volume.master.*.id, count.index)
  }

  console {
    type        = "pty"
    target_port = 0
  }

  cpu = {
    mode = "host-passthrough"
  }

  network_interface {
    network_id = libvirt_network.net.id
    hostname   = "${var.cluster_id}-master-${count.index}.${var.cluster_domain}"
    addresses  = [var.libvirt_master_ips[count.index]]
  }
}

data "libvirt_network_dns_host_template" "bootstrap" {
  count    = var.bootstrap_dns ? 1 : 0
  ip       = var.libvirt_bootstrap_ip
  hostname = "api.${var.cluster_domain}"
}

data "libvirt_network_dns_host_template" "masters" {
  count    = var.master_count
  ip       = var.libvirt_master_ips[count.index]
  hostname = "api.${var.cluster_domain}"
}

data "libvirt_network_dns_host_template" "bootstrap_int" {
  count    = var.bootstrap_dns ? 1 : 0
  ip       = var.libvirt_bootstrap_ip
  hostname = "api-int.${var.cluster_domain}"
}

data "libvirt_network_dns_host_template" "masters_int" {
  count    = var.master_count
  ip       = var.libvirt_master_ips[count.index]
  hostname = "api-int.${var.cluster_domain}"
}

data "libvirt_network_dns_host_template" "etcds" {
  count    = var.master_count
  ip       = var.libvirt_master_ips[count.index]
  hostname = "etcd-${count.index}.${var.cluster_domain}"
}

data "libvirt_network_dns_srv_template" "etcd_cluster" {
  count    = var.master_count
  service  = "etcd-server-ssl"
  protocol = "tcp"
  domain   = var.cluster_domain
  port     = 2380
  weight   = 10
  target   = "etcd-${count.index}.${var.cluster_domain}"
}

