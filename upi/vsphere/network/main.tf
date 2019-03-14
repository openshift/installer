data "external" "ping" {
  program = ["bash", "${path.root}/network/cidr_to_ip.sh"]

  query = {
    cidr = "${var.machine_cidr}"
  }
}

output "ip_list" "usable" {
  value = "${data.external.ping.result.ipaddress}"
}
