resource "openstack_compute_instance_v2" "etcd_node" {
  count           = "${var.tectonic_etcd_count}"
  name            = "${var.tectonic_cluster_name}_etcd_node_${count.index}"
  image_id        = "${var.tectonic_openstack_image_id}"
  flavor_id       = "${var.tectonic_openstack_flavor_id}"
  key_pair        = "${openstack_compute_keypair_v2.k8s_keypair.name}"
  security_groups = ["${openstack_compute_secgroup_v2.etcd_group.name}"]

  metadata {
    role = "etcd"
  }

  user_data    = "${ignition_config.etcd.*.rendered[count.index]}"
  config_drive = false

  network {
    uuid = "${openstack_networking_network_v2.network.id}"
  }
}

resource "openstack_compute_secgroup_v2" "etcd_group" {
  name        = "${var.tectonic_cluster_name}_etcd_group"
  description = "security group for etcd: SSH and etcd client / cluster"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 2379
    to_port     = 2380
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}
