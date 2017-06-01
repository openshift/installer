resource "null_resource" "bootstrap" {
  # Without depends_on, this remote-exec may start before the kubeconfig copy.  # Terraform only does one task at a time, so it would try to bootstrap  # Kubernetes and Tectonic while no Kubelets are running. Ensure all nodes  # receive a kubeconfig before proceeding with bootkube and tectonic.  #depends_on = ["null_resource.kubeconfig-masters"]

  connection {
    type    = "ssh"
    host    = "${module.masters.ip_address[0]}"
    user    = "core"
    timeout = "60m"
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
      "sudo systemctl start tectonic",
    ]
  }
}
