resource "aws_route53_record" "tnc_internal" {
  zone_id = "${local.private_zone_id}"
  name    = "${var.tectonic_cluster_name}-tnc.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${local.tnc_elb_dns_name}"
    zone_id                = "${local.tnc_elb_zone_id}"
    evaluate_target_health = true
  }
}
