resource "aws_lb" "api_internal" {
  name                             = "${var.cluster_name}-int"
  load_balancer_type               = "network"
  subnets                          = ["${local.master_subnet_ids}"]
  internal                         = true
  enable_cross_zone_load_balancing = true
  idle_timeout                     = 3600

  tags = "${merge(map(
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

resource "aws_lb" "api_external" {
  name                             = "${var.cluster_name}-ext"
  load_balancer_type               = "network"
  subnets                          = ["${local.master_subnet_ids}"]
  internal                         = false
  enable_cross_zone_load_balancing = true
  idle_timeout                     = 3600

  tags = "${merge(map(
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

resource "aws_lb_target_group" "api_internal" {
  name     = "${var.cluster_name}-api-int"
  protocol = "TCP"
  port     = 6443
  vpc_id   = "${local.vpc_id}"

  target_type = "ip"

  tags = "${merge(map(
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"

  health_check {
    healthy_threshold   = 3
    unhealthy_threshold = 3
    interval            = 10
    port                = 6443
    protocol            = "TCP"
  }
}

resource "aws_lb_target_group" "api_external" {
  name     = "${var.cluster_name}-api-ext"
  protocol = "TCP"
  port     = 6443
  vpc_id   = "${local.vpc_id}"

  target_type = "ip"

  tags = "${merge(map(
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"

  health_check {
    healthy_threshold   = 3
    unhealthy_threshold = 3
    interval            = 10
    port                = 6443
    protocol            = "TCP"
  }
}

resource "aws_lb_target_group" "services" {
  name     = "${var.cluster_name}-services"
  protocol = "TCP"
  port     = 49500
  vpc_id   = "${local.vpc_id}"

  target_type = "ip"

  tags = "${merge(map(
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"

  health_check {
    healthy_threshold   = 3
    unhealthy_threshold = 3
    interval            = 10
    port                = 49500
    protocol            = "TCP"
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
  port              = "49500"

  default_action {
    target_group_arn = "${aws_lb_target_group.services.arn}"
    type             = "forward"
  }
}

resource "aws_lb_listener" "api_external_api" {
  count = "${var.public_master_endpoints ? 1 : 0}"

  load_balancer_arn = "${aws_lb.api_external.arn}"
  protocol          = "TCP"
  port              = "6443"

  default_action {
    target_group_arn = "${aws_lb_target_group.api_external.arn}"
    type             = "forward"
  }
}
