resource "aws_autoscaling_group" "workers" {
  name                 = "${var.cluster_name}-worker"
  desired_capacity     = "${var.worker_count}"
  max_size             = "${var.worker_count * 3}"
  min_size             = "${var.worker_count}"
  launch_configuration = "${aws_launch_configuration.worker_conf.id}"
  vpc_zone_identifier  = ["${aws_subnet.worker_subnet.*.id}"]
  tag                  = "key=Name,value=${var.cluster_name}-worker"
  propagate_at_launch  = true

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_launch_configuration" "worker_conf" {
  instance_type   = "${var.worker_ec2_type}"
  image_id        = "${data.aws_ami.coreos_ami.image_id}"
  name_prefix     = "${var.cluster_name}-worker-"
  key_name        = "${aws_key_pair.ssh-key.key_name}"
  security_groups = ["${aws_security_group.worker_sec_group.id}"]

  user_data = "${ignition_config.worker.rendered}"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "worker_sec_group" {
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
