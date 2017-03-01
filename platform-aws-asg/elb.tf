resource "aws_elb" "api-internal" {
  name            = "${var.cluster_name}-api-internal"
  subnets         = ["${aws_subnet.master_subnet.*.id}"]
  internal        = true
  security_groups = ["${aws_security_group.master_sec_group.id}"]

  listener {
    instance_port     = 443
    instance_protocol = "tcp"
    lb_port           = 443
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 10255
    instance_protocol = "tcp"
    lb_port           = 10255
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:10255/healthz"
    interval            = 30
  }

  tags {
    Name              = "${var.cluster_name}-api-internal"
    KubernetesCluster = "${var.cluster_name}"
  }
}

resource "aws_elb" "api-external" {
  name            = "${var.cluster_name}-api-external"
  subnets         = ["${aws_subnet.az_subnet_pub.*.id}"]
  internal        = false
  security_groups = ["${aws_security_group.master_sec_group.id}"]

  listener {
    instance_port     = 22
    instance_protocol = "tcp"
    lb_port           = 22
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 443
    instance_protocol = "tcp"
    lb_port           = 443
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "TCP:22"
    interval            = 30
  }

  tags {
    Name              = "${var.cluster_name}-api-external"
    KubernetesCluster = "${var.cluster_name}"
  }
}

resource "aws_elb" "console" {
  name            = "${var.cluster_name}-console"
  subnets         = ["${aws_subnet.az_subnet_pub.*.id}"]
  internal        = false
  security_groups = ["${aws_security_group.master_sec_group.id}"]

  listener {
    instance_port     = 32001
    instance_protocol = "tcp"
    lb_port           = 80
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 32000
    instance_protocol = "tcp"
    lb_port           = 443
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:32002/healthz"
    interval            = 5
  }

  tags {
    Name              = "${var.cluster_name}-console"
    KubernetesCluster = "${var.cluster_name}"
  }
}
