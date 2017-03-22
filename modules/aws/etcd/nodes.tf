data "aws_ami" "coreos_ami" {
  most_recent = true

  filter {
    name   = "name"
    values = ["CoreOS-${var.tectonic_cl_channel}-*"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "owner-id"
    values = ["595879546273"]
  }
}

resource "aws_instance" "etcd_node" {
  count                  = "${length(var.external_endpoints) == 0 ? var.node_count : 0}"
  ami                    = "${data.aws_ami.coreos_ami.image_id}"
  instance_type          = "t2.medium"
  subnet_id              = "${var.etcd_subnets[count.index]}"
  key_name               = "${var.ssh_key}"
  user_data              = "${ignition_config.etcd.*.rendered[count.index]}"
  vpc_security_group_ids = ["${aws_security_group.etcd_sec_group.id}"]

  tags {
    Name              = "${var.tectonic_cluster_name}-etcd-${count.index}"
    KubernetesCluster = "${var.tectonic_cluster_name}"
  }
}
