data "aws_route53_zone" "base" {
  name = "${var.base_domain}"
}

locals {
  public_zone_id = "${data.aws_route53_zone.base.zone_id}"

  zone_id = "${var.private_zone_id}"
}

resource "aws_route53_record" "api_external" {
  zone_id = "${local.public_zone_id}"
  name    = "${var.cluster_name}-api.${var.base_domain}"
  type    = "A"

  alias {
    name                   = "${var.api_external_lb_dns_name}"
    zone_id                = "${var.api_external_lb_zone_id}"
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "api_internal" {
  zone_id = "${var.private_zone_id}"
  name    = "${var.cluster_name}-api.${var.base_domain}"
  type    = "A"

  alias {
    name                   = "${var.api_internal_lb_dns_name}"
    zone_id                = "${var.api_internal_lb_zone_id}"
    evaluate_target_health = true
  }
}
