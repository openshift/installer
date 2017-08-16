resource "aws_route53_record" "worker_nodes" {
  count   = "${var.worker_count}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-worker-${count.index}"
  type    = "A"
  ttl     = "60"
  records = ["${var.worker_ips[count.index]}"]
}
