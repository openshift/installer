resource "aws_route53_zone" "tectonic-int" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"
  name   = "${var.tectonic_domain}"
}

data "aws_route53_zone" "tectonic-ext" {
  name = "${var.tectonic_domain}"
}
