data "aws_route53_zone" "tectonic" {
  name = "${var.base_domain}"
}

resource "aws_route53_record" "tectonic_api" {
  count   = "1"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-k8s"
  type    = "A"
  ttl     = "60"
  records = ["${var.api_ips}"]
}

resource "aws_route53_record" "tectonic-console" {
  count   = "${var.tectonic_vanilla_k8s ? 0 : 1}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}"
  type    = "A"
  ttl     = "60"
  records = ["${var.worker_ips}"]
}
