resource "null_resource" "kubeconfig" {
  count = "${length(var.tectonic_controller_domains) + length(var.tectonic_worker_domains)}"
  depends_on = ["module.tectonic"]

  connection {
    type = "ssh"
    host = "${element(concat(var.tectonic_controller_domains, var.tectonic_worker_domains), count.index)}"
    user = "core"
  }

  provisioner "file" {
    source = "${path.cwd}/generated/kubeconfig"
    destination = "$HOME/kubeconfig"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mv /home/core/kubeconfig /etc/kubernetes/kubeconfig",
    ]
  }
}

resource "null_resource" "bootstrap" {
  depends_on = ["module.tectonic"]

  connection {
    type = "ssh"
    host = "${element(var.tectonic_controller_domains, 0)}"
    user = "core"
  }

  provisioner "file" {
    source = "${path.cwd}/generated"
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
