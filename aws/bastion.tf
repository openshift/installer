/*
    The bastion instance is used only for debugging during development.
    This file should be removed once the templates are stable and working correctly.
*/
resource "aws_security_group" "bastion_sec_group" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 22
    to_port     = 22
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "ignition_config" "bastion" {
  files = [
    "${ignition_file.etcd-endpoints.id}",
  ]
}

resource "aws_instance" "bastion_node" {
  ami                         = "${data.aws_ami.coreos_ami.image_id}"
  instance_type               = "t2.small"
  subnet_id                   = "${aws_subnet.az_subnet_pub.0.id}"
  key_name                    = "${aws_key_pair.ssh-key.key_name}"
  vpc_security_group_ids      = ["${aws_security_group.bastion_sec_group.id}"]
  user_data                   = "${ignition_config.bastion.rendered}"
  associate_public_ip_address = true

  tags {
    Name = "bastion"
  }
}
