resource "aws_route53_record" "etcd_srv_discover" {
  count   = "${var.dns_enabled ? 1 : 0}"
  name    = "_etcd-server._tcp"
  type    = "SRV"
  zone_id = "${var.dns_zone_id}"
  records = ["${formatlist("0 0 2380 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "300"
}

resource "aws_route53_record" "etcd_srv_client" {
  count   = "${var.dns_enabled ? 1 : 0}"
  name    = "_etcd-client._tcp"
  type    = "SRV"
  zone_id = "${var.dns_zone_id}"
  records = ["${formatlist("0 0 2379 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "60"
}

resource "aws_route53_record" "etc_a_nodes" {
  count   = "${var.dns_enabled ? var.instance_count : 0}"
  type    = "A"
  ttl     = "60"
  zone_id = "${var.dns_zone_id}"
  name    = "${var.cluster_name}-etcd-${count.index}"
  records = ["${aws_instance.etcd_node.*.private_ip[count.index]}"]
}
