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
  records = ["${openstack_networking_port_v2.master.*.all_fixed_ips[count.index]}"]
}

resource "aws_route53_record" "worker_nodes" {
  count   = "${var.tectonic_worker_count}"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
  name    = "${var.tectonic_cluster_name}-worker-${count.index}"
  type    = "A"
  ttl     = "60"
  records = ["${openstack_networking_port_v2.worker.*.all_fixed_ips[count.index]}"]
}

# etcd

resource "aws_route53_record" "etcd_srv_discover" {
  count = "${var.tectonic_experimental ? 0 : 1}"

  name    = "${var.tectonic_etcd_tls_enabled ? "_etcd-server-ssl._tcp" : "_etcd-server._tcp"}"
  type    = "SRV"
  records = ["${formatlist("0 0 2380 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "300"
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
}

resource "aws_route53_record" "etcd_srv_client" {
  count = "${var.tectonic_experimental ? 0 : 1}"

  name    = "${var.tectonic_etcd_tls_enabled ? "_etcd-client-ssl._tcp" : "_etcd-client._tcp"}"
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
  records = ["${openstack_networking_port_v2.etcd.*.all_fixed_ips[count.index]}"]
  zone_id = "${data.aws_route53_zone.tectonic.zone_id}"
}
