locals {
  network      = "${cidrhost(var.machine_cidr,0)}"
  ip_addresses = ["${coalescelist(var.ip_addresses, data.template_file.ip_address.*.rendered)}"]
}

data "external" "ip_address" {
  count = "${length(var.ip_addresses) == 0 ? var.instance_count : 0}"

  program = ["bash", "${path.module}/cidr_to_ip.sh"]

  query = {
    hostname   = "${var.name}-${count.index}.${var.cluster_domain}"
    ipam       = "${var.ipam}"
    ipam_token = "${var.ipam_token}"
  }

  depends_on = ["null_resource.ip_address"]
}

data "template_file" "ip_address" {
  count = "${length(var.ip_addresses) == 0 ? var.instance_count : 0}"

  template = "${lookup(data.external.ip_address.*.result[count.index], "ip_address")}"
}

resource "null_resource" "ip_address" {
  count = "${length(var.ip_addresses) == 0 ? var.instance_count : 0}"

  provisioner "local-exec" {
    command = <<EOF
echo '{"network":"${local.network}","hostname":"${var.name}-${count.index}.${var.cluster_domain}","ipam":"${var.ipam}","ipam_token":"${var.ipam_token}"}' | ${path.module}/cidr_to_ip.sh
EOF
  }

  provisioner "local-exec" {
    when = "destroy"

    command = <<EOF
curl -s "http://${var.ipam}/api/removeHost.php?apiapp=address&apitoken=${var.ipam_token}&host=${var.name}-${count.index}.${var.cluster_domain}"
EOF
  }
}
