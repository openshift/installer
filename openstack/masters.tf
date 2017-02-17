resource "openstack_compute_instance_v2" "control_node" {
  count           = "${var.controller_count}"
  name            = "control_node_${count.index}"
  image_id        = "${var.image_id}"
  flavor_id       = "${var.flavor_id}"
  key_pair        = "${openstack_compute_keypair_v2.k8s_keypair.name}"
  security_groups = ["${openstack_compute_secgroup_v2.k8s_control_group.name}"]

  metadata {
    role = "controller"
  }

  user_data    = "${data.template_file.userdata-master.rendered}"
  config_drive = false

  # connection {
  #   user        = "core"
  #   private_key = "${tls_private_key.core.private_key_pem}"
  # }
  # # copy something so we wait until the host is ready
  # provisioner "file" {
  #   source      = "../kubelet.master"
  #   destination = "/home/core/kubelet.master"
  # }
  # provisioner "remote-exec" {
  #   inline = [
  #     "sudo mv /home/core/kubelet.master /etc/systemd/system/kubelet.service",
  #     "chmod +x ./init-master.sh",
  #     "sudo ./init-master.sh local",
  #   ]
  # }
}

resource "openstack_compute_secgroup_v2" "k8s_control_group" {
  name        = "k8s_control_group"
  description = "security group for k8s controllers: SSH and https"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 443
    to_port     = 443
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
