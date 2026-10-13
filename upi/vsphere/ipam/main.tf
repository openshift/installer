locals {
  network      = cidrhost(var.machine_cidr, 0)
  hostnames    = length(var.static_ip_addresses) == 0 ? var.hostnames : []
  ip_addresses = length(var.static_ip_addresses) == 0 ? [for result in null_resource.ip_address : jsondecode(data.http.getip[result.triggers.hostname].body)[result.triggers.hostname]] : var.static_ip_addresses
}

data "http" "getip" {
  for_each = null_resource.ip_address

  url = "http://${var.ipam}/api/getIPs.php?apiapp=address&apitoken=${var.ipam_token}&domain=${null_resource.ip_address[each.key].triggers.hostname}"

  request_headers = {
    Accept = "application/json"
  }
}

resource "null_resource" "ip_address" {
  for_each = local.hostnames

  triggers = {
    ipam       = var.ipam
    ipam_token = var.ipam_token
    network    = local.network
    hostname   = each.key
  }

  provisioner "local-exec" {
    command = <<EOF
echo '{"network":"${self.triggers.network}","hostname":"${self.triggers.hostname}","ipam":"${self.triggers.ipam}","ipam_token":"${self.triggers.ipam_token}"}' | ${path.module}/cidr_to_ip.sh
EOF

  }
  provisioner "local-exec" {
    when = destroy

    command = <<EOF
curl -s "http://${self.triggers.ipam}/api/removeHost.php?apiapp=address&apitoken=${self.triggers.ipam_token}&host=${self.triggers.hostname}"
EOF

  }
}
