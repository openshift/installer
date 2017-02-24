data "template_file" "userdata" {
  count    = "${var.node_count}"
  template = "${file("${path.module}/userdata.yml")}"

  vars {
    node_name   = "etcd-${count.index}.${var.tectonic_domain}"
    etcd_domain = "${var.tectonic_domain}"
  }
}

resource "aws_instance" "etcd_node" {
  count                  = "${var.node_count}"
  ami                    = "${var.coreos_ami}"
  instance_type          = "t2.medium"
  subnet_id              = "${var.etcd_subnets[count.index]}"
  key_name               = "${var.ssh_key}"
  user_data              = "${data.template_file.userdata.*.rendered[count.index]}"
  vpc_security_group_ids = ["${aws_security_group.etcd_sec_group.id}"]

  tags {
    Name = "${var.cluster_name}-etcd-${count.index}"
  }
}
