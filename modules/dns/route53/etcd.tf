resource "aws_route53_record" "etcd_srv_discover" {
  count   = "1"
  name    = "${var.etcd_tls_enabled ? "_etcd-server-ssl._tcp" : "_etcd-server._tcp"}"
  type    = "SRV"
  zone_id = "${aws_route53_zone.tectonic.zone_id}"
  records = ["${formatlist("0 0 2380 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "300"
}

resource "aws_route53_record" "etcd_srv_client" {
  count   = "1"
  name    = "${var.etcd_tls_enabled ? "_etcd-client-ssl._tcp" : "_etcd-client._tcp"}"
  type    = "SRV"
  zone_id = "${aws_route53_zone.tectonic.zone_id}"
  records = ["${formatlist("0 0 2379 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "60"
}

resource "aws_route53_record" "etc_a_nodes" {
  count   = "${var.etcd_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-etcd-${count.index}"
  records = ["${var.etcd_ips[count.index]}"]
}
