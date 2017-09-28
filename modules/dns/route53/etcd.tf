resource "aws_route53_record" "etcd_a_nodes" {
  count   = "${var.tectonic_experimental ? 0 : var.etcd_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-etcd-${count.index}"
  records = ["${var.etcd_ips[count.index]}"]
}
