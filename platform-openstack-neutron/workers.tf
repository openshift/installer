resource "openstack_compute_instance_v2" "worker_node" {
  count     = "${var.tectonic_worker_count}"
  name      = "${var.tectonic_cluster_name}_worker_node_${count.index}"
  image_id  = "${var.tectonic_openstack_image_id}"
  flavor_id = "${var.tectonic_openstack_flavor_id}"
  key_pair  = "${openstack_compute_keypair_v2.k8s_keypair.name}"

  metadata {
    role = "worker"
  }

  network {
    floating_ip = "${openstack_compute_floatingip_v2.worker.*.address[count.index]}"
    uuid = "${openstack_networking_network_v2.network.id}"
  }

  user_data    = "${ignition_config.worker.*.rendered[count.index]}"
  config_drive = false
}
