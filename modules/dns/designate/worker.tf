resource "openstack_dns_recordset_v2" "worker_nodes" {
  count   = "${var.worker_count}"
  zone_id = "${openstack_dns_zone_v2.tectonic.id}"
  name    = "${var.cluster_name}-worker-${count.index}.${var.base_domain}."
  type    = "A"
  ttl     = "60"
  records = ["${var.worker_ips[count.index]}"]
}
