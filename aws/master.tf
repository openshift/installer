resource "aws_autoscaling_group" "masters" {
  name                 = "${var.cluster_name}-masters"
  desired_capacity     = "${var.master_count}"
  max_size             = "${var.master_count * 3}"
  min_size             = "${var.master_count}"
  launch_configuration = "${aws_launch_configuration.master_conf.id}"
  vpc_zone_identifier  = ["${aws_subnet.master_subnet.*.id}"]

  load_balancers = ["${aws_elb.api-internal.id}", "${aws_elb.api-external.id}"]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_launch_configuration" "master_conf" {
  instance_type   = "${var.master_ec2_type}"
  image_id        = "${data.aws_ami.coreos_ami.image_id}"
  name_prefix     = "${var.cluster_name}-master-"
  key_name        = "${aws_key_pair.ssh-key.key_name}"
  security_groups = ["${aws_security_group.master_sec_group.id}"]

  user_data = "${ignition_config.master.rendered}"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "master_sec_group" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 22
    to_port     = 22
  }

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 443
    to_port     = 443
  }

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 10255
    to_port     = 10255
  }

  ingress {
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 10250
    to_port     = 10250
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
  }
}
