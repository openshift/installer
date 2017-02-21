resource "aws_route53_zone" etcd_zone {
  vpc_id = "${data.aws_vpc.etcd_vpc.id}"
  name   = "${var.etcd_domain}"
}

resource "aws_route53_record" "etcd_srv_discover" {
  name    = "_etcd-server._tcp"
  type    = "SRV"
  zone_id = "${aws_route53_zone.etcd_zone.id}"
  records = ["${formatlist("0 0 2380 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "300"
}

resource "aws_route53_record" "etcd_srv_client" {
  name    = "_etcd-client._tcp"
  type    = "SRV"
  zone_id = "${aws_route53_zone.etcd_zone.id}"
  records = ["${formatlist("0 0 2379 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "60"
}

resource "aws_route53_record" "etc_a_nodes" {
  count   = "${var.node_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.etcd_zone.id}"
  name    = "node-${count.index}"
  records = ["${aws_instance.etcd_node.*.private_ip[count.index]}"]
}
