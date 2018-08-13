provider "aws" {
  region  = "${var.tectonic_aws_region}"
  profile = "${var.tectonic_aws_profile}"
  version = "1.8.0"

  assume_role {
    role_arn     = "${var.tectonic_aws_installer_role == "" ? "" : "${var.tectonic_aws_installer_role}"}"
    session_name = "TECTONIC_INSTALLER_${var.tectonic_cluster_name}"
  }
}

resource "aws_route53_record" "tectonic_tnc_a" {
  depends_on = ["aws_route53_record.tectonic_tnc_cname"]
  count      = "1"
  zone_id    = "${local.private_zone_id}"
  name       = "${var.tectonic_cluster_name}-tnc.${var.tectonic_base_domain}"
  type       = "A"

  alias {
    name                   = "${local.tnc_elb_dns_name}"
    zone_id                = "${local.tnc_elb_zone_id}"
    evaluate_target_health = true
  }
}
