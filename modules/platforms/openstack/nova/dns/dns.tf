data "aws_route53_zone" "tectonic" {
  name = "${var.base_domain}"
}

resource "aws_route53_record" "tectonic-api" {
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-k8s"
  type    = "A"
  ttl     = "60"
  records = ["${var.tectonic_api_records}"]
}

resource "aws_route53_record" "tectonic-console" {
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}"
  type    = "A"
  ttl     = "60"
  records = ["${var.tectonic_console_records}"]
}

resource "aws_route53_record" "etcd" {
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-etc"
  type    = "A"
  ttl     = "60"
  records = ["${var.etcd_records}"]
}

resource "aws_route53_record" "master_nodes" {
  count   = "${var.master_count}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-master-${count.index}"
  type    = "A"
  ttl     = "60"
  records = ["${var.master_records[count.index]}"]
}

resource "aws_route53_record" "worker_nodes" {
  count   = "${var.worker_count}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-worker-${count.index}"
  type    = "A"
  ttl     = "60"
  records = ["${var.worker_records[count.index]}"]
}
