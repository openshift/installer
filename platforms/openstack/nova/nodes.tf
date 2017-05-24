resource "openstack_compute_instance_v2" "master_node" {
  count     = "${var.tectonic_master_count}"
  name      = "${var.tectonic_cluster_name}_master_node_${count.index}"
  image_id  = "${var.tectonic_openstack_image_id}"
  flavor_id = "${var.tectonic_openstack_flavor_id}"

  security_groups = [
    "${module.master_nodes.secgroup_master_name}",
    "${module.master_nodes.secgroup_self_hosted_etcd_name}",
  ]

  network {
    name = "${var.tectonic_openstack_network_name}"
  }

  metadata {
    role = "master"
  }

  user_data    = "${module.master_nodes.user_data[count.index]}"
  config_drive = false
}

resource "openstack_compute_instance_v2" "worker_node" {
  count           = "${var.tectonic_worker_count}"
  name            = "${var.tectonic_cluster_name}_worker_node_${count.index}"
  image_id        = "${var.tectonic_openstack_image_id}"
  flavor_id       = "${var.tectonic_openstack_flavor_id}"
  security_groups = ["${module.worker_nodes.secgroup_node_name}"]

  network {
    name = "${var.tectonic_openstack_network_name}"
  }

  metadata {
    role = "worker"
  }

  user_data    = "${module.worker_nodes.user_data[count.index]}"
  config_drive = false
}

resource "openstack_compute_instance_v2" "etcd_node" {
  count           = "${var.tectonic_experimental ? 0 : var.tectonic_etcd_count}"
  name            = "${var.tectonic_cluster_name}_etcd_node_${count.index}"
  image_id        = "${var.tectonic_openstack_image_id}"
  flavor_id       = "${var.tectonic_openstack_flavor_id}"
  security_groups = ["${module.etcd.secgroup_name}"]

  network {
    name = "${var.tectonic_openstack_network_name}"
  }

  metadata {
    role = "etcd"
  }

  user_data    = "${module.etcd.user_data[count.index]}"
  config_drive = false
}

resource "null_resource" "tectonic" {
  depends_on = ["module.tectonic"]

  connection {
    host        = "${openstack_compute_instance_v2.master_node.*.access_ip_v4[0]}"
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
      "sudo systemctl start ${var.tectonic_vanilla_k8s ? "bootkube.service" : "tectonic.service"}",
    ]
  }
}
