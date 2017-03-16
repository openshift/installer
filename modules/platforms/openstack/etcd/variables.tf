// The amount of etcd nodes to be created.
// Example: `3`
variable "count" {
  type = "string"
}

// The name of the cluster.
// The etcd hostnames will be prefixed with this.
variable "cluster_name" {
  type = "string"
}

// The public keys for the core user.
variable core_public_keys {
  type = "list"
}

output "user_data" {
  value = ["${ignition_config.etcd.*.rendered}"]
}

output "secgroup_name" {
  value = "${openstack_compute_secgroup_v2.etcd.name}"
}
