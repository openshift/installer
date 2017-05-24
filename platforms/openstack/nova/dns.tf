# tectonic

data "aws_route53_zone" "tectonic" {
  name = "${var.tectonic_base_domain}"
}

resource "aws_route53_record" "tectonic-api" {
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.tectonic_cluster_name}-k8s"
  type    = "A"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.master_node.*.access_ip_v4}"]
}

resource "aws_route53_record" "tectonic-console" {
  count   = "${var.tectonic_vanilla_k8s ? 0 : 1}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.tectonic_cluster_name}"
  type    = "A"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.worker_node.*.access_ip_v4}"]
}

# master/worker

resource "aws_route53_record" "master_nodes" {
  count   = "${var.tectonic_master_count}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.tectonic_cluster_name}-master-${count.index}"
  type    = "A"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.master_node.*.access_ip_v4[count.index]}"]
}

resource "aws_route53_record" "worker_nodes" {
  count   = "${var.tectonic_worker_count}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.tectonic_cluster_name}-worker-${count.index}"
  type    = "A"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.worker_node.*.access_ip_v4[count.index]}"]
}

# etcd

resource "aws_route53_record" "etcd_srv_discover" {
  name    = "_etcd-server._tcp"
  count   = "${var.tectonic_experimental ? 0 : 1}"
  type    = "SRV"
  records = ["${formatlist("0 0 2380 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "300"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
}

resource "aws_route53_record" "etcd_srv_client" {
  name    = "_etcd-client._tcp"
  count   = "${var.tectonic_experimental ? 0 : 1}"
  type    = "SRV"
  records = ["${formatlist("0 0 2379 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "60"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
}

resource "aws_route53_record" "etc_a_nodes" {
  count   = "${var.tectonic_experimental ? 0 : var.tectonic_etcd_count}"
  type    = "A"
  ttl     = "60"
  name    = "${var.tectonic_cluster_name}-etcd-${count.index}"
  records = ["${openstack_compute_instance_v2.etcd_node.*.access_ip_v4[count.index]}"]
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
}
