output "secgroup_master_names" {
  value = ["${compact(list(
    openstack_compute_secgroup_v2.default.name,
    openstack_compute_secgroup_v2.master.name,
    var.tectonic_experimental ? openstack_compute_secgroup_v2.self_hosted_etcd.name : "",
  ))}"]
}

output "secgroup_master_ids" {
  value = ["${compact(list(
    openstack_compute_secgroup_v2.default.id,
    openstack_compute_secgroup_v2.master.id,
    var.tectonic_experimental ? openstack_compute_secgroup_v2.self_hosted_etcd.id : "",
  ))}"]
}

output "secgroup_node_names" {
  value = ["${openstack_compute_secgroup_v2.node.name}", "${openstack_compute_secgroup_v2.default.name}"]
}

output "secgroup_node_ids" {
  value = ["${openstack_compute_secgroup_v2.node.id}", "${openstack_compute_secgroup_v2.default.id}"]
}

output "secgroup_etcd_names" {
  value = ["${openstack_compute_secgroup_v2.etcd.name}", "${openstack_compute_secgroup_v2.default.name}"]
}

output "secgroup_etcd_ids" {
  value = ["${openstack_compute_secgroup_v2.etcd.id}", "${openstack_compute_secgroup_v2.default.id}"]
}
