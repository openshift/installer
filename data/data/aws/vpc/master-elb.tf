resource "aws_lb" "api_internal" {
  name                             = "${var.cluster_id}-int"
  load_balancer_type               = "network"
  subnets                          = local.private_subnet_ids
  internal                         = true
  enable_cross_zone_load_balancing = true
  idle_timeout                     = 3600

  tags = merge(
    {
      "Name" = "${var.cluster_id}-int"
    },
    var.tags,
  )

  depends_on = [aws_internet_gateway.igw]
}

resource "aws_lb" "api_external" {
  name                             = "${var.cluster_id}-ext"
  load_balancer_type               = "network"
  subnets                          = local.public_subnet_ids
  internal                         = false
  enable_cross_zone_load_balancing = true
  idle_timeout                     = 3600

  tags = merge(
    {
      "Name" = "${var.cluster_id}-ext"
    },
    var.tags,
  )

  depends_on = [aws_internet_gateway.igw]
}

resource "aws_lb_target_group" "api_internal" {
  name     = "${var.cluster_id}-aint"
  protocol = "TCP"
  port     = 6443
  vpc_id   = local.vpc_id

  target_type = "ip"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-aint"
    },
    var.tags,
  )

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
  vpc_id   = local.vpc_id

  target_type = "ip"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-aext"
    },
    var.tags,
  )

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    interval            = 10
    port                = 6443
    protocol            = "HTTPS"
    path                = "/readyz"
  }
}

resource "aws_lb_target_group" "services" {
  name     = "${var.cluster_id}-sint"
  protocol = "TCP"
  port     = 22623
  vpc_id   = local.vpc_id

  target_type = "ip"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-sint"
    },
    var.tags,
  )

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
  load_balancer_arn = aws_lb.api_internal.arn
  protocol          = "TCP"
  port              = "6443"

  default_action {
    target_group_arn = aws_lb_target_group.api_internal.arn
    type             = "forward"
  }
}

resource "aws_lb_listener" "api_internal_services" {
  load_balancer_arn = aws_lb.api_internal.arn
  protocol          = "TCP"
  port              = "22623"

  default_action {
    target_group_arn = aws_lb_target_group.services.arn
    type             = "forward"
  }
}

resource "aws_lb_listener" "api_external_api" {
  count = var.public_master_endpoints ? 1 : 0

  load_balancer_arn = aws_lb.api_external.arn
  protocol          = "TCP"
  port              = "6443"

  default_action {
    target_group_arn = aws_lb_target_group.api_external.arn
    type             = "forward"
  }
}

