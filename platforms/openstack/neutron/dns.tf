# tectonic

data "aws_route53_zone" "tectonic" {
  name = "${var.tectonic_base_domain}"
}

resource "aws_route53_record" "tectonic-api" {
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.tectonic_cluster_name}-k8s"
  type    = "A"
  ttl     = "60"
  records = ["${openstack_networking_floatingip_v2.master.*.address}"]
}

resource "aws_route53_record" "tectonic-console" {
  count   = "${var.tectonic_vanilla_k8s ? 0 : 1}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.tectonic_cluster_name}"
  type    = "A"
  ttl     = "60"
  records = ["${openstack_networking_floatingip_v2.worker.*.address}"]
}

# master/worker

resource "aws_route53_record" "master_nodes" {
  count   = "${var.tectonic_master_count}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.tectonic_cluster_name}-master-${count.index}"
  type    = "A"
  ttl     = "60"
  records = ["${element(openstack_networking_port_v2.master.*.fixed_ip.0.ip_address, count.index)}"]
}

resource "aws_route53_record" "worker_nodes" {
  count   = "${var.tectonic_worker_count}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.tectonic_cluster_name}-worker-${count.index}"
  type    = "A"
  ttl     = "60"
  records = ["${element(openstack_networking_port_v2.worker.*.fixed_ip.0.ip_address, count.index)}"]
}

# etcd

resource "aws_route53_record" "etcd_srv_discover" {
  name    = "_etcd-server._tcp"
  type    = "SRV"
  records = ["${formatlist("0 0 2380 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "300"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
}

resource "aws_route53_record" "etcd_srv_client" {
  name    = "_etcd-client._tcp"
  type    = "SRV"
  records = ["${formatlist("0 0 2379 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "60"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
}

resource "aws_route53_record" "etc_a_nodes" {
  count   = "${var.tectonic_etcd_count}"
  type    = "A"
  ttl     = "60"
  name    = "${var.tectonic_cluster_name}-etcd-${count.index}"
  records = ["${openstack_compute_instance_v2.etcd_node.*.access_ip_v4[count.index]}"]
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
}
