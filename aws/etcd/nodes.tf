data "aws_availability_zones" "zones" {}

data "aws_ami" "coreos_ami" {
  most_recent = true

  filter {
    name   = "name"
    values = ["CoreOS-stable-*"]
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
  count         = "${var.node_count}"
  ami           = "${data.aws_ami.coreos_ami.id}"
  instance_type = "t2.medium"
  subnet_id     = "${data.aws_subnet.az_subnet.*.id[count.index]}"
  key_name      = "${aws_key_pair.ssh-key.id}"
  user_data     = "${data.template_file.userdata.*.rendered[count.index]}"

  tags {
    Name = "node-${count.index}"
  }
}
