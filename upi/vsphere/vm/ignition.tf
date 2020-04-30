locals {
  ignition_encoded = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition)}"
}

data "ignition_file" "hostname" {
  for_each = var.hostnames_ip_addresses

  filesystem = "root"
  path       = "/etc/hostname"
  mode       = "420"

  content {
    content = element(split(".", each.key), 0)
  }
}

data "ignition_file" "static_ip" {
  for_each = var.hostnames_ip_addresses

  filesystem = "root"
  path       = "/etc/sysconfig/network-scripts/ifcfg-ens192"
  mode       = "420"

  content {
    content = templatefile("${path.module}/ifcfg.tmpl", {
      dns_addresses = var.dns_addresses,
      machine_cidr  = var.machine_cidr
      //ip_address     = var.hostnames_ip_addresses[count.index].value
      ip_address     = each.value
      cluster_domain = var.cluster_domain
    })
  }
}

data "ignition_config" "ign" {
  for_each = var.hostnames_ip_addresses

  append {
    source = local.ignition_encoded
  }

  files = [
    data.ignition_file.hostname[each.key].rendered,
    data.ignition_file.static_ip[each.key].rendered,
  ]
}

