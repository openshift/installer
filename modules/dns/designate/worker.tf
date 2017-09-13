resource "openstack_dns_recordset_v2" "worker_nodes" {
  count   = "${var.worker_count}"
  zone_id = "${data.openstack_dns_zone_v2.tectonic.id}"
  name    = "${var.cluster_name}-worker-${count.index}.${var.base_domain}."
  type    = "A"
  ttl     = "60"
  records = ["${var.worker_ips[count.index]}"]
}

resource "openstack_dns_recordset_v2" "worker_nodes_public" {
  count   = "${var.worker_count}"
  zone_id = "${data.openstack_dns_zone_v2.tectonic.id}"
  name    = "${var.cluster_name}-worker-${count.index}-public.${var.base_domain}."
  type    = "A"
  ttl     = "60"
  records = ["${var.worker_public_ips[count.index]}"]
}
