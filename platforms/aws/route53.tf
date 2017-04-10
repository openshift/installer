data "aws_route53_zone" "tectonic-ext" {
  count = "${var.tectonic_aws_external_vpc_public}"
  name  = "${var.tectonic_base_domain}"
}

resource "aws_route53_zone" "tectonic-int" {
  vpc_id = "${module.vpc.vpc_id}"
  name   = "${var.tectonic_base_domain}"
}
