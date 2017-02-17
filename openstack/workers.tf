resource "openstack_compute_instance_v2" "worker_node" {
  count     = "${var.worker_count}"
  name      = "worker_node_${count.index}"
  image_id  = "${var.image_id}"
  flavor_id = "${var.flavor_id}"
  key_pair  = "${openstack_compute_keypair_v2.k8s_keypair.name}"

  metadata {
    role = "worker"
  }

  user_data    = "${data.template_file.userdata-worker.rendered}"
  config_drive = false

  # connection {
  #   user        = "core"
  #   private_key = "${tls_private_key.core.private_key_pem}"
  # }
  # provisioner "file" {
  #   source      = "../kubelet.master"
  #   destination = "/home/core/kubelet.worker"
  # }
}
