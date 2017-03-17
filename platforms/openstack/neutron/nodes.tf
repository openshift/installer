resource "openstack_compute_instance_v2" "etcd_node" {
  count           = "${var.tectonic_etcd_count}"
  name            = "${var.tectonic_cluster_name}_etcd_node_${count.index}"
  image_id        = "${var.tectonic_openstack_image_id}"
  flavor_id       = "${var.tectonic_openstack_flavor_id}"
  security_groups = ["${module.etcd.secgroup_name}"]

  metadata {
    role = "etcd"
  }

  network {
    uuid = "${module.network.id}"
  }

  user_data    = "${module.etcd.user_data[count.index]}"
  config_drive = false
}

resource "openstack_compute_instance_v2" "master_node" {
  count           = "${var.tectonic_master_count}"
  name            = "${var.tectonic_cluster_name}_master_node_${count.index}"
  image_id        = "${var.tectonic_openstack_image_id}"
  flavor_id       = "${var.tectonic_openstack_flavor_id}"
  security_groups = ["${module.nodes.master_secgroup_name}"]

  metadata {
    role = "master"
  }

  network {
    floating_ip = "${module.network.master_floating_ips[count.index]}"
    uuid        = "${module.network.id}"
  }

  user_data    = "${module.nodes.master_user_data[count.index]}"
  config_drive = false
}

resource "openstack_compute_instance_v2" "worker_node" {
  count     = "${var.tectonic_worker_count}"
  name      = "${var.tectonic_cluster_name}_worker_node_${count.index}"
  image_id  = "${var.tectonic_openstack_image_id}"
  flavor_id = "${var.tectonic_openstack_flavor_id}"

  metadata {
    role = "worker"
  }

  network {
    floating_ip = "${module.network.worker_floating_ips[count.index]}"
    uuid        = "${module.network.id}"
  }

  user_data    = "${module.nodes.worker_user_data[count.index]}"
  config_drive = false
}

resource "null_resource" "tectonic" {
  depends_on = ["module.tectonic"]

  connection {
    host        = "${module.network.worker_floating_ips[0]}"
    private_key = "${module.secrets.core_private_key_pem}"
    user        = "core"
  }

  provisioner "file" {
    source      = "${path.cwd}/generated"
    destination = "$HOME/tectonic"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mkdir -p /opt",
      "sudo rm -rf /opt/tectonic",
      "sudo mv /home/core/tectonic /opt/",
      "sudo systemctl start tectonic",
    ]
  }
}
