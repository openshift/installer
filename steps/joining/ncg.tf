resource "aws_route53_record" "ngc_internal" {
  zone_id = "${local.private_zone_id}"
  name    = "${var.tectonic_cluster_name}-ncg.${var.tectonic_base_domain}"
  type    = "A"

  alias {
    name                   = "${local.ncg_elb_dns_name}"
    zone_id                = "${local.ncg_elb_zone_id}"
    evaluate_target_health = true
  }
}
