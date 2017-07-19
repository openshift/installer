module "designate" {
  source = "../../../modules/dns/designate"

  cluster_name = "${var.tectonic_cluster_name}"
  base_domain  = "${var.tectonic_base_domain}"
  admin_email  = "${var.tectonic_admin_email}"

  master_count = "${var.tectonic_master_count}"
  master_ips   = "${openstack_networking_floatingip_v2.master.*.address}"

  worker_count = "${var.tectonic_worker_count}"
  worker_ips   = "${openstack_networking_floatingip_v2.worker.*.address}"

  etcd_count       = "${var.tectonic_experimental ? 0 : var.tectonic_etcd_count}"
  etcd_ips         = "${openstack_networking_port_v2.etcd.*.all_fixed_ips}"
  etcd_tls_enabled = "${var.tectonic_etcd_tls_enabled}"
}
