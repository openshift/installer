# resource "null_resource" "bootkube" {
#   triggers {
#     master-nodes = "${join(",",aws_elb.api-external.instances)}"
#   }
#   connection {
#     host  = "${aws_elb.api-external.dns_name}"
#     user  = "core"
#     agent = true
#   }
#   provisioner "file" {
#     source      = "${path.cwd}/assets"
#     destination = "$HOME/assets"
#   }
#   provisioner "remote-exec" {
#     inline = [
#       "sudo mv /home/core/assets /opt/bootkube/",
#       "sudo chmod a+x /opt/bootkube/assets/bootkube-start",
#       "sudo systemctl start bootkube",
#     ]
#   }
# }

