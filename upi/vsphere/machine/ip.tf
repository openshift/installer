locals {
  ip_addresses = ["${data.external.ip_address.*.result.ip_address}"]
  ips_exist    = "${length(compact(local.ip_addresses)) == var.instance_count}"
}

data "external" "ip_address" {
  count = "${var.instance_count}"

  program = ["bash", "${path.module}/cidr_to_ip.sh"]

  query = {
    hostname   = "${var.name}-${count.index}.${var.cluster_domain}"
    ipam       = "${var.ipam}"
    ipam_token = "${var.ipam_token}"
  }
}

resource "null_resource" "ip_address" {
  count = "${var.instance_count}"

  provisioner "local-exec" {
    command = <<EOF
echo '{"cidr":"${var.machine_cidr}","hostname":"${var.name}-${count.index}.${var.cluster_domain}","ipam":"${var.ipam}","ipam_token":"${var.ipam_token}"}' | ${path.module}/cidr_to_ip.sh
EOF
  }

  provisioner "local-exec" {
    when = "destroy"

    command = <<EOF
curl -s "http://${var.ipam}/api/removeHost.php?apiapp=address&apitoken=${var.ipam_token}&host=${var.name}-${count.index}.${var.cluster_domain}"
EOF
  }
}
