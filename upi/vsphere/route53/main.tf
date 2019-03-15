data "aws_route53_zone" "base" {
  name = "${var.base_domain}"
}

resource "aws_route53_zone" "cluster" {
  name          = "${var.cluster_domain}"
  force_destroy = true
}

resource "aws_route53_record" "api_external" {
  type    = "A"
  ttl     = "60"
  zone_id = "${data.aws_route53_zone.base.zone_id}"
  name    = "api.${var.cluster_domain}"

  weighted_routing_policy {
    weight = 90
  }

  set_identifier = "api_external"
  records        = ["${var.etcd_ip_addresses}"]
}

resource "aws_route53_record" "api_internal" {
  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "api.${var.cluster_domain}"

  weighted_routing_policy {
    weight = 90
  }

  set_identifier = "api_internal"
  records        = ["${var.etcd_ip_addresses}"]
}

resource "aws_route53_record" "workers" {
  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "workers.${var.cluster_domain}"

  weighted_routing_policy {
    weight = 90
  }

  set_identifier = "workers"
  records        = ["${var.worker_ips}"]
}

resource "aws_route53_record" "masters" {
  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "masters.${var.cluster_domain}"

  weighted_routing_policy {
    weight = 90
  }

  set_identifier = "masters"
  records        = ["${var.etcd_ip_addresses}"]
}

resource "aws_route53_record" "bootstrap" {
  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "bootstrap.${var.cluster_domain}"
  records = ["${var.bootstrap_ip}"]
}

resource "aws_route53_record" "etcd_a_nodes" {
  count   = "${var.etcd_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "etcd-${count.index}.${var.cluster_domain}"
  records = ["${var.etcd_ip_addresses[count.index]}"]
}

resource "aws_route53_record" "etcd_cluster" {
  type    = "SRV"
  ttl     = "60"
  zone_id = "${aws_route53_zone.cluster.zone_id}"
  name    = "_etcd-server-ssl._tcp"
  records = ["${formatlist("0 10 2380 %s", aws_route53_record.etcd_a_nodes.*.fqdn)}"]
}
