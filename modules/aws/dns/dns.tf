data "aws_route53_zone" "tectonic-ext" {
  name = "${var.tectonic_base_domain}"
}

data "aws_vpc" "cluster_vpc" {
  id = "${var.vpc_id}"
}

resource "aws_route53_zone" "tectonic-int" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"
  name   = "${var.tectonic_base_domain}"
}

resource "aws_route53_record" "api-internal" {
  zone_id = "${aws_route53_zone.tectonic-int.zone_id}"
  name    = "${var.tectonic_dns_name}-k8s.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${var.api-internal-elb["dns_name"]}"
    zone_id                = "${var.api-internal-elb["zone_id"]}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "api-external" {
  zone_id = "${data.aws_route53_zone.tectonic-ext.zone_id}"
  name    = "${var.tectonic_dns_name}-k8s.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${var.api-external-elb["dns_name"]}"
    zone_id                = "${var.api-external-elb["zone_id"]}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "ingress-public" {
  zone_id = "${data.aws_route53_zone.tectonic-ext.zone_id}"
  name    = "${var.tectonic_dns_name}.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${var.console-elb["dns_name"]}"
    zone_id                = "${var.console-elb["zone_id"]}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "ingress-private" {
  zone_id = "${aws_route53_zone.tectonic-int.zone_id}"
  name    = "${var.tectonic_dns_name}.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${var.console-elb["dns_name"]}"
    zone_id                = "${var.console-elb["zone_id"]}"
    evaluate_target_health = true
  }
}
