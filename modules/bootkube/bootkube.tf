resource "null_resource" "copy_assets" {
  # Changes to any instance of the cluster requires re-provisioning
  triggers {
    cluster_instance_ids = "${join(" ", var.trigger_ids)}"
  }

  # Bootstrap script can run on any instance of the cluster
  # So we just choose the first in this case
  connection {
    user        = "core"
    private_key = "${var.core_private_key}"
    host        = "${var.hosts[0]}"
  }

  provisioner "file" {
    source      = "${var.assets_dir}"
    destination = "/home/core/assets"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo mv /home/core/assets /opt/bootkube/",
      "sudo chmod a+x /opt/bootkube/assets/bootkube-start",
      "sudo systemctl start bootkube",
    ]
  }
}
