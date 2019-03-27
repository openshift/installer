data "aws_route53_zone" "base" {
  name = "${var.base_domain}"
}

resource "aws_route53_zone" "cluster" {
  name          = "${var.cluster_domain}"
  force_destroy = true

  tags = "${map(
    "Name", "${var.cluster_domain}",
	"Platform", "vSphere")}"
}

resource "aws_route53_record" "name_server" {
  name    = "${var.cluster_domain}"
  type    = "NS"
  ttl     = "300"
  zone_id = "${data.aws_route53_zone.base.zone_id}"
  records = ["${aws_route53_zone.cluster.name_servers}"]
}

resource "aws_route53_record" "api" {
  type           = "A"
  ttl            = "60"
  zone_id        = "${aws_route53_zone.cluster.zone_id}"
  name           = "api.${var.cluster_domain}"
  set_identifier = "api"
  records        = ["${compact(concat(list(var.bootstrap_ip), var.control_plane_ips))}"]

  weighted_routing_policy {
    weight = 90
  }
}

resource "aws_route53_record" "etcd_a_nodes" {
  count = "${length(var.control_plane_ips)}"

  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "etcd-${count.index}.${var.cluster_domain}"
  records = ["${var.control_plane_ips[count.index]}"]
}

resource "aws_route53_record" "etcd_cluster" {
  type    = "SRV"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "_etcd-server-ssl._tcp"
  records = ["${formatlist("0 10 2380 %s", aws_route53_record.etcd_a_nodes.*.fqdn)}"]
}

resource "aws_route53_record" "ingress" {
  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "*.apps.${var.cluster_domain}"
  records = ["${var.compute_ips}"]
}

resource "aws_route53_record" "control_plane_nodes" {
  count = "${length(var.control_plane_ips)}"

  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "control-plane-${count.index}.${var.cluster_domain}"
  records = ["${var.control_plane_ips[count.index]}"]
}

resource "aws_route53_record" "compute_nodes" {
  count = "${length(var.compute_ips)}"

  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "compute-${count.index}.${var.cluster_domain}"
  records = ["${var.compute_ips[count.index]}"]
}
