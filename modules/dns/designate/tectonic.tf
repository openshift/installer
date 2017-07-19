resource "openstack_dns_zone_v2" "tectonic" {
  count = "1"
  name  = "${var.base_domain}."
  email = "${var.admin_email}"
  ttl   = "60"
}

resource "openstack_dns_recordset_v2" "tectonic-api" {
  count   = "1"
  zone_id = "${openstack_dns_zone_v2.tectonic.id}"
  name    = "${var.cluster_name}-k8s.${var.base_domain}."
  type    = "A"
  ttl     = "60"
  records = ["${var.master_ips}"]
}

resource "openstack_dns_recordset_v2" "tectonic-console" {
  count   = "1"
  zone_id = "${openstack_dns_zone_v2.tectonic.id}"
  name    = "${var.cluster_name}.${var.base_domain}."
  type    = "A"
  ttl     = "60"
  records = ["${var.worker_ips}"]
}
