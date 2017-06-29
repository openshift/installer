resource "null_resource" "etcd_secrets" {
  count = "${var.tectonic_experimental ? 0 : var.tectonic_etcd_count }"

  connection {
    type    = "ssh"
    host    = "${element(module.etcd.ip_address, count.index)}"
    user    = "core"
    timeout = "60m"
  }

  provisioner "file" {
    content     = "${module.bootkube.etcd_ca_crt_pem}"
    destination = "$HOME/etcd_ca.crt"
  }

  provisioner "file" {
    content     = "${module.bootkube.etcd_client_crt_pem}"
    destination = "$HOME/etcd_client.crt"
  }

  provisioner "file" {
    content     = "${module.bootkube.etcd_client_key_pem}"
    destination = "$HOME/etcd_client.key"
  }

  provisioner "file" {
    content     = "${module.bootkube.etcd_peer_crt_pem}"
    destination = "$HOME/etcd_peer.crt"
  }

  provisioner "file" {
    content     = "${module.bootkube.etcd_peer_key_pem}"
    destination = "$HOME/etcd_peer.key"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mkdir -p /etc/ssl/etcd",
      "sudo mv /home/core/etcd_ca.crt /etc/ssl/etcd/ca.crt",
      "sudo mv /home/core/etcd_client.crt /etc/ssl/etcd/client.crt",
      "sudo mv /home/core/etcd_client.key /etc/ssl/etcd/client.key",
      "sudo mv /home/core/etcd_peer.key /etc/ssl/etcd/peer.key",
      "sudo mv /home/core/etcd_peer.crt /etc/ssl/etcd/peer.crt",
    ]
  }
}

resource "null_resource" "bootstrap" {
  # Without depends_on, this remote-exec may start before the kubeconfig copy.  # Terraform only does one task at a time, so it would try to bootstrap  # Kubernetes and Tectonic while no Kubelets are running. Ensure all nodes  # receive a kubeconfig before proceeding with bootkube and tectonic.  #depends_on = ["null_resource.kubeconfig-masters"]

  connection {
    type        = "ssh"
    host        = "${module.masters.ip_address[0]}"
    user        = "core"
    timeout     = "60m"
    private_key = "${file(var.tectonic_vmware_ssh_private_key_path != "" ? pathexpand(var.tectonic_vmware_ssh_private_key_path) : "/dev/null")}"
  }

  provisioner "file" {
    source      = "./generated"
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
