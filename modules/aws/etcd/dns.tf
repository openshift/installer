resource "aws_route53_record" "etc_a_nodes" {
  count   = "${var.dns_enabled ? var.instance_count : 0}"
  type    = "A"
  ttl     = "60"
  zone_id = "${var.dns_zone_id}"
  name    = "${var.cluster_name}-etcd-${count.index}"
  records = ["${aws_instance.etcd_node.*.private_ip[count.index]}"]
}
