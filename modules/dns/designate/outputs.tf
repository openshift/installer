output "etcd_a_nodes" {
  value = "${openstack_dns_recordset_v2.etcd_a_nodes.*.records}"
}

output "master_nodes" {
  value = "${openstack_dns_recordset_v2.master_nodes.*.records}"
}

output "worker_nodes" {
  value = "${openstack_dns_recordset_v2.worker_nodes.*.records}"
}

output "worker_nodes_public" {
  value = "${openstack_dns_recordset_v2.worker_nodes_public.*.records}"
}
