resource "null_resource" "bootstrapper" {
  triggers {
    endpoint = "${var.bootstrapping_host}"
  }

  connection {
    host  = "${var.bootstrapping_host}"
    user  = "core"
    agent = true
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
    ]
  }
}
