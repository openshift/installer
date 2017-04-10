data "aws_route53_zone" "tectonic-ext" {
  name = "${var.tectonic_base_domain}"
}

resource "aws_route53_zone" "tectonic-int" {
  vpc_id = "${module.vpc.vpc_id}"
  name   = "${var.tectonic_base_domain}"
}
