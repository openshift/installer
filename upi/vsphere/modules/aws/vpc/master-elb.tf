resource "aws_lb" "api_internal" {
  name                             = "${var.cluster_id}-int"
  load_balancer_type               = "network"
  subnets                          = "${var.private_subnet_id}"
  internal                         = true
  enable_cross_zone_load_balancing = true
  idle_timeout                     = 3600

  tags = "${merge(map(
    "Name", "${var.cluster_id}-int",
  ), var.tags)}"

  timeouts {
    create = "20m"
  }
}

resource "aws_lb" "api_external" {
  name                             = "${var.cluster_id}-ext"
  load_balancer_type               = "network"
  subnets                          = "${var.public_subnet_id}"
  internal                         = false
  enable_cross_zone_load_balancing = true
  idle_timeout                     = 3600

  tags = "${merge(map(
    "Name", "${var.cluster_id}-ext",
  ), var.tags)}"

  timeouts {
    create = "20m"
  }
}

resource "aws_lb_target_group" "api_internal" {
  name     = "${var.cluster_id}-aint"
  protocol = "TCP"
  port     = 6443
  vpc_id   = "${var.vpc_id}"

  target_type = "ip"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-aint",
  ), var.tags)}"

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    interval            = 10
    port                = 6443
    protocol            = "HTTPS"
    path                = "/readyz"
  }
}

resource "aws_lb_target_group" "api_external" {
  name     = "${var.cluster_id}-aext"
  protocol = "TCP"
  port     = 6443
  vpc_id   = "${var.vpc_id}"

  target_type = "ip"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-aext",
  ), var.tags)}"

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    interval            = 10
    port                = 6443
    protocol            = "HTTPS"
    path                = "/readyz"
  }
}

resource "aws_lb_target_group" "ingress_external_https" {
  name     = "${var.cluster_id}-iext-https"
  protocol = "TCP"
  port     = 443
  vpc_id   = "${var.vpc_id}"

  target_type = "ip"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-iext-https",
  ), var.tags)}"
}

resource "aws_lb_target_group" "ingress_external_http" {
  name     = "${var.cluster_id}-iext-http"
  protocol = "TCP"
  port     = 80
  vpc_id   = "${var.vpc_id}"

  target_type = "ip"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-iext-http",
  ), var.tags)}"
}

resource "aws_lb_target_group" "services" {
  name     = "${var.cluster_id}-sint"
  protocol = "TCP"
  port     = 22623
  vpc_id   = "${var.vpc_id}"

  target_type = "ip"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-sint",
  ), var.tags)}"

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    interval            = 10
    port                = 22623
    protocol            = "HTTPS"
    path                = "/healthz"
  }
}

resource "aws_lb_listener" "api_internal_api" {
  load_balancer_arn = "${aws_lb.api_internal.arn}"
  protocol          = "TCP"
  port              = "6443"

  default_action {
    target_group_arn = "${aws_lb_target_group.api_internal.arn}"
    type             = "forward"
  }
}

resource "aws_lb_listener" "api_internal_services" {
  load_balancer_arn = "${aws_lb.api_internal.arn}"
  protocol          = "TCP"
  port              = "22623"

  default_action {
    target_group_arn = "${aws_lb_target_group.services.arn}"
    type             = "forward"
  }
}

resource "aws_lb_listener" "ingress_external_ingress_https" {
  load_balancer_arn = "${aws_lb.api_external.arn}"
  protocol          = "TCP"
  port              = "443"

  default_action {
    target_group_arn = "${aws_lb_target_group.ingress_external_https.arn}"
    type             = "forward"
  }
}

resource "aws_lb_listener" "ingress_external_ingress_http" {
  load_balancer_arn = "${aws_lb.api_external.arn}"
  protocol          = "TCP"
  port              = "80"

  default_action {
    target_group_arn = "${aws_lb_target_group.ingress_external_http.arn}"
    type             = "forward"
  }
}

resource "aws_lb_listener" "api_external_api" {
  count = "${var.public_control_plane_endpoints ? 1 : 0}"

  load_balancer_arn = "${aws_lb.api_external.arn}"
  protocol          = "TCP"
  port              = "6443"

  default_action {
    target_group_arn = "${aws_lb_target_group.api_external.arn}"
    type             = "forward"
  }
}
