resource "aws_route53_record" "etcd_srv_discover" {
  name    = "_etcd-server._tcp"
  type    = "SRV"
  zone_id = "${var.dns_zone}"
  records = ["${formatlist("0 0 2380 %s", aws_route53_record.etcd_a_nodes.*.fqdn)}"]
  ttl     = "300"
}

resource "aws_route53_record" "etcd_srv_client" {
  name    = "_etcd-client._tcp"
  type    = "SRV"
  zone_id = "${var.dns_zone}"
  records = ["${formatlist("0 0 2379 %s", aws_route53_record.etcd_a_nodes.*.fqdn)}"]
  ttl     = "60"
}

resource "aws_route53_record" "etcd_a_nodes" {
  count   = "${var.node_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${var.dns_zone}"
  name    = "etcd-${count.index}"
  records = ["${aws_instance.etcd_node.*.private_ip[count.index]}"]
}
