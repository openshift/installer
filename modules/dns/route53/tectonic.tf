resource "aws_route53_zone" "tectonic" {
  count = "1"
  name  = "${var.base_domain}"
}

resource "aws_route53_record" "tectonic_api" {
  count   = "1"
  zone_id = "${aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-k8s"
  type    = "A"
  ttl     = "60"
  records = ["${var.master_ips}"]
}

resource "aws_route53_record" "tectonic_console" {
  count   = "1"
  zone_id = "${aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}"
  type    = "A"
  ttl     = "60"
  records = ["${var.worker_ips}"]
}
