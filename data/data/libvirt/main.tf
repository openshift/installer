provider "libvirt" {
  uri = var.libvirt_uri
}

module "volume" {
  source = "./volume"

  cluster_id = var.cluster_id
  image      = var.os_image
}

module "bootstrap" {
  source = "./bootstrap"

  addresses      = [var.libvirt_bootstrap_ip]
  base_volume_id = module.volume.coreos_base_volume_id
  cluster_id     = var.cluster_id
  ignition       = var.ignition_bootstrap
  network_id     = libvirt_network.net.id
}

resource "libvirt_volume" "master" {
  count          = var.master_count
  name           = "${var.cluster_id}-master-${count.index}"
  base_volume_id = module.volume.coreos_base_volume_id
}

resource "libvirt_ignition" "master" {
  name    = "${var.cluster_id}-master.ign"
  content = var.ignition_master
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
      for_each = [flatten(
        [
          data.libvirt_network_dns_srv_template.etcd_cluster.*.rendered,
        ],
      )]
      content {
        # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
        # which keys might be set in maps assigned here, so it has
        # produced a comprehensive set here. Consider simplifying
        # this after confirming which keys can be set in practice.

        domain   = lookup(srvs.value, "domain", null)
        port     = lookup(srvs.value, "port", null)
        priority = lookup(srvs.value, "priority", null)
        protocol = lookup(srvs.value, "protocol", null)
        service  = lookup(srvs.value, "service", null)
        target   = lookup(srvs.value, "target", null)
        weight   = lookup(srvs.value, "weight", null)
      }
    }

    dynamic "hosts" {
      for_each = [flatten(
        [
          data.libvirt_network_dns_host_template.bootstrap.*.rendered,
          data.libvirt_network_dns_host_template.bootstrap_int.*.rendered,
          data.libvirt_network_dns_host_template.masters.*.rendered,
          data.libvirt_network_dns_host_template.masters_int.*.rendered,
          data.libvirt_network_dns_host_template.etcds.*.rendered,
        ],
      )]
      content {
        # TF-UPGRADE-TODO: The automatic upgrade tool can't predict
        # which keys might be set in maps assigned here, so it has
        # produced a comprehensive set here. Consider simplifying
        # this after confirming which keys can be set in practice.

        hostname = lookup(hosts.value, "hostname", null)
        ip       = lookup(hosts.value, "ip", null)
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
    hostname   = "${var.cluster_id}-master-${count.index}"
    # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
    # force an interpolation expression to be interpreted as a list by wrapping it
    # in an extra set of list brackets. That form was supported for compatibilty in
    # v0.11, but is no longer supported in Terraform v0.12.
    #
    # If the expression in the following list itself returns a list, remove the
    # brackets to avoid interpretation as a list of lists. If the expression
    # returns a single list item then leave it as-is and remove this TODO comment.
    addresses = [var.libvirt_master_ips[count.index]]
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

