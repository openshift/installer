output "etc_a_nodes" {
  value = "${openstack_dns_recordset_v2.etc_a_nodes.*.name}"
}
