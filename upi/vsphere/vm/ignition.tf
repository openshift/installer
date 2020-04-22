locals {
  ignition_encoded = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition)}"
}

data "ignition_file" "hostname" {
  for_each = var.hostnames_ip_addresses

  path = "/etc/hostname"
  mode = "420"

  content {
    content = element(split(".", each.key), 0)
  }
}

data "ignition_file" "static_ip" {
  for_each = var.hostnames_ip_addresses

  path = "/etc/NetworkManager/system-connections/ens192"
  mode = "384"

  content {
    content = templatefile("${path.module}/nm-keyfile.tmpl", {
      dns_addresses = var.dns_addresses,
      machine_cidr  = var.machine_cidr
      //ip_address    = var.hostnames_ip_addresses[count.index].value
      ip_address    = each.value
    })
  }
}

data "ignition_config" "ign" {
  for_each = var.hostnames_ip_addresses

  merge {
    source = local.ignition_encoded
  }

  files = [
    data.ignition_file.hostname[each.key].rendered,
    data.ignition_file.static_ip[each.key].rendered,
  ]
}

