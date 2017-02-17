data "aws_route53_zone" "tectonic" {
  name = "${var.base_domain}"
}

resource "aws_route53_record" "tectonic-api" {
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-k8s"
  type    = "A"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.control_node.*.access_ip_v4}"]
}

# resource "aws_route53_record" "tectonic-console" {
#   zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
#   name    = "${var.cluster_name}"
#   type    = "A"
#   ttl     = "60"
#   records = [""]
# }

resource "aws_route53_record" "etcd" {
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.cluster_name}-etc"
  type    = "A"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.etcd_node.*.access_ip_v4}"]
}
