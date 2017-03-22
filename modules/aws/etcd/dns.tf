resource "aws_route53_record" "etcd_srv_discover" {
  count   = "${length(var.external_endpoints) == 0 ? 1 : 0}"
  name    = "_etcd-server._tcp"
  type    = "SRV"
  zone_id = "${var.dns_zone}"
  records = ["${formatlist("0 0 2380 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "300"
}

resource "aws_route53_record" "etcd_srv_client" {
  count   = "${length(var.external_endpoints) == 0 ? 1 : 0}"
  name    = "_etcd-client._tcp"
  type    = "SRV"
  zone_id = "${var.dns_zone}"
  records = ["${formatlist("0 0 2379 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "60"
}

resource "aws_route53_record" "etc_a_nodes" {
  count   = "${length(var.external_endpoints) == 0 ? var.node_count : 0}"
  type    = "A"
  ttl     = "60"
  zone_id = "${var.dns_zone}"
  name    = "${var.tectonic_cluster_name}-etcd-${count.index}"
  records = ["${aws_instance.etcd_node.*.private_ip[count.index]}"]
}
