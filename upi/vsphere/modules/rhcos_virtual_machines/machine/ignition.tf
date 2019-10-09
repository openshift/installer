provider "ignition" {
  version = "1.1.0"
}

locals {
  mask                   = "${element(split("/", var.machine_cidr), 1)}"
  gw                     = "${cidrhost(var.machine_cidr,1)}"
  dns_addresses_rendered = "${join("\n", formatlist("%s=%s", null_resource.generate_nm_dns_addresses.*.triggers.dns_item, var.dns_addresses))}"

  ignition_encoded = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition)}"
}

/* generate_nm_dns_addresses just creates
 * DNS1 and DNS2 to be
 * used in local dns_addresses_rendered
 */
resource "null_resource" "generate_nm_dns_addresses" {
  count = "${length(var.dns_addresses)}"

  triggers {
    dns_item = "DNS${count.index + 1}"
  }
}

data "ignition_file" "hostname" {
  count = "${var.instance_count}"

  filesystem = "root"
  path       = "/etc/hostname"
  mode       = "420"

  content {
    content = "${var.name}-${count.index}"
  }
}

data "ignition_file" "static_ip" {
  count = "${var.instance_count}"

  filesystem = "root"
  path       = "/etc/sysconfig/network-scripts/ifcfg-ens192"
  mode       = "420"

  content {
    content = <<EOF
TYPE=Ethernet
BOOTPROTO=none
NAME=ens192
DEVICE=ens192
ONBOOT=yes
IPADDR=${var.ip_addresses[count.index]}
PREFIX=${local.mask}
GATEWAY=${local.gw}
DOMAIN=${var.cluster_domain}
${local.dns_addresses_rendered}
EOF
  }
}

data "ignition_config" "ign" {
  count = "${var.instance_count}"

  append {
    source = "${local.ignition_encoded}"
  }

  files = [
    "${data.ignition_file.hostname.*.id[count.index]}",
    "${data.ignition_file.static_ip.*.id[count.index]}",
  ]
}
