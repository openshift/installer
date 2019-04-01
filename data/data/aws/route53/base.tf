data "aws_route53_zone" "public" {
  name = "${var.base_domain}"
}

resource "aws_route53_zone" "int" {
  name          = "${var.cluster_domain}"
  force_destroy = true

  vpc {
    vpc_id = "${var.vpc_id}"
  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}-int",
  ), var.tags)}"

  depends_on = ["aws_route53_record.api_external"]
}

resource "aws_route53_record" "api_external" {
  zone_id = "${data.aws_route53_zone.public.zone_id}"
  name    = "api.${var.cluster_domain}"
  type    = "A"

  alias {
    name                   = "${var.api_external_lb_dns_name}"
    zone_id                = "${var.api_external_lb_zone_id}"
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api_internal" {
  zone_id = "${aws_route53_zone.int.zone_id}"
  name    = "api.${var.cluster_domain}"
  type    = "A"

  alias {
    name                   = "${var.api_internal_lb_dns_name}"
    zone_id                = "${var.api_internal_lb_zone_id}"
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "etcd_a_nodes" {
  count   = "${var.etcd_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.int.zone_id}"
  name    = "etcd-${count.index}.${var.cluster_domain}"
  records = ["${var.etcd_ip_addresses[count.index]}"]
}

resource "aws_route53_record" "etcd_cluster" {
  type    = "SRV"
  ttl     = "60"
  zone_id = "${aws_route53_zone.int.zone_id}"
  name    = "_etcd-server-ssl._tcp"
  records = ["${formatlist("0 10 2380 %s", aws_route53_record.etcd_a_nodes.*.fqdn)}"]
}
