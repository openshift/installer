resource "aws_route53_zone" "tectonic-int" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"
  name   = "${var.tectonic_base_domain}"
}

data "aws_route53_zone" "tectonic-ext" {
  name = "${var.tectonic_base_domain}"
}

resource "aws_route53_record" "api-internal" {
  zone_id = "${aws_route53_zone.tectonic-int.zone_id}"
  name    = "${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${aws_elb.api-internal.dns_name}"
    zone_id                = "${aws_elb.api-internal.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "api-external" {
  zone_id = "${data.aws_route53_zone.tectonic-ext.zone_id}"
  name    = "${var.tectonic_cluster_name}-k8s.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${aws_elb.api-external.dns_name}"
    zone_id                = "${aws_elb.api-external.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "ingress-public" {
  zone_id = "${data.aws_route53_zone.tectonic-ext.zone_id}"
  name    = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${aws_elb.console.dns_name}"
    zone_id                = "${aws_elb.console.zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "ingress-private" {
  zone_id = "${aws_route53_zone.tectonic-int.zone_id}"
  name    = "${var.tectonic_cluster_name}.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${aws_elb.console.dns_name}"
    zone_id                = "${aws_elb.console.zone_id}"
    evaluate_target_health = true
  }
}
