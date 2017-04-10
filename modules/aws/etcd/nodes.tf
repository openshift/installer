data "aws_ami" "coreos_ami" {
  most_recent = true

  filter {
    name   = "name"
    values = ["CoreOS-${var.cl_channel}-*"]
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
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"
  ami   = "${data.aws_ami.coreos_ami.image_id}"

  instance_type          = "${var.ec2_type}"
  subnet_id              = "${var.subnets[count.index % var.az_count]}"
  key_name               = "${var.ssh_key}"
  user_data              = "${ignition_config.etcd.*.rendered[count.index]}"
  vpc_security_group_ids = ["${aws_security_group.etcd_sec_group.id}"]

  tags {
    Name              = "${var.cluster_name}-etcd-${count.index}"
    KubernetesCluster = "${var.cluster_name}"
  }
}
