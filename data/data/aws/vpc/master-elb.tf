# AWS NLBs (Network Load Balancers, or the aws_lb resource) do not support
# IPv6.  Instead, we must use class load balancers (the aws_elb resource).

resource "aws_lb" "api_internal" {
  name                             = "${var.cluster_id}-int"
  load_balancer_type               = "network"
  subnets                          = data.aws_subnet.private.*.id
  internal                         = true
  enable_cross_zone_load_balancing = true

  tags = merge(
    {
      "Name" = "${var.cluster_id}-int"
    },
    var.tags,
  )

  timeouts {
    create = "20m"
  }

  depends_on = [aws_internet_gateway.igw]

  count = var.use_ipv6 == false ? 1 : 0
}

resource "aws_lb" "api_external" {
  count = local.public_endpoints && var.use_ipv6 == false ? 1 : 0

  name                             = "${var.cluster_id}-ext"
  load_balancer_type               = "network"
  subnets                          = data.aws_subnet.public.*.id
  internal                         = false
  enable_cross_zone_load_balancing = true

  tags = merge(
    {
      "Name" = "${var.cluster_id}-ext"
    },
    var.tags,
  )

  timeouts {
    create = "20m"
  }

  depends_on = [aws_internet_gateway.igw]
}

resource "aws_lb_target_group" "api_internal" {
  name     = "${var.cluster_id}-aint"
  protocol = "TCP"
  port     = 6443
  vpc_id   = data.aws_vpc.cluster_vpc.id

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

  count = var.use_ipv6 == false ? 1 : 0
}

resource "aws_lb_target_group" "api_external" {
  count = local.public_endpoints && var.use_ipv6 == false ? 1 : 0

  name     = "${var.cluster_id}-aext"
  protocol = "TCP"
  port     = 6443
  vpc_id   = data.aws_vpc.cluster_vpc.id

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
  vpc_id   = data.aws_vpc.cluster_vpc.id

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

  count = var.use_ipv6 == false ? 1 : 0
}

resource "aws_lb_listener" "api_internal_api" {
  load_balancer_arn = aws_lb.api_internal[0].arn
  protocol          = "TCP"
  port              = "6443"

  default_action {
    target_group_arn = aws_lb_target_group.api_internal[0].arn
    type             = "forward"
  }

  count = var.use_ipv6 == false ? 1 : 0
}

resource "aws_lb_listener" "api_internal_services" {
  load_balancer_arn = aws_lb.api_internal[0].arn
  protocol          = "TCP"
  port              = "22623"

  default_action {
    target_group_arn = aws_lb_target_group.services[0].arn
    type             = "forward"
  }

  count = var.use_ipv6 == false ? 1 : 0
}

resource "aws_lb_listener" "api_external_api" {
  count = local.public_endpoints && var.use_ipv6 == false ? 1 : 0

  load_balancer_arn = aws_lb.api_external[0].arn
  protocol          = "TCP"
  port              = "6443"

  default_action {
    target_group_arn = aws_lb_target_group.api_external[0].arn
    type             = "forward"
  }
}

#################################
### Begin IPv6 load balancers ###
#################################

resource "aws_security_group" "api" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  timeouts {
    create = "20m"
  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}_api_sg",
  ), var.tags)}"

  count = var.use_ipv6 == true ? 1 : 0
}

resource "aws_security_group_rule" "api_ingress_api_v6" {
  type              = "ingress"
  security_group_id = "${aws_security_group.api[0].id}"

  protocol         = "tcp"
  ipv6_cidr_blocks = ["::/0"]
  from_port        = 6443
  to_port          = 6443

  count = var.use_ipv6 == true ? 1 : 0
}

resource "aws_security_group_rule" "api_ingress_mcs_v6" {
  type              = "ingress"
  security_group_id = "${aws_security_group.api[0].id}"

  protocol         = "tcp"
  ipv6_cidr_blocks = ["::/0"]
  from_port        = 22623
  to_port          = 22623

  count = var.use_ipv6 == true ? 1 : 0
}

resource "aws_elb" "api_internal" {
  name            = "${var.cluster_id}-int"
  subnets         = data.aws_subnet.private.*.id
  internal        = true
  security_groups = ["${aws_security_group.master.id}"]

  idle_timeout                = 3600
  connection_draining         = true
  connection_draining_timeout = 300

  listener {
    instance_port     = 6443
    instance_protocol = "tcp"
    lb_port           = 6443
    lb_protocol       = "tcp"
  }

  listener {
    instance_port     = 22623
    instance_protocol = "tcp"
    lb_port           = 22623
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "SSL:6443"
    interval            = 5
  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}-int",
  ), var.tags)}"

  count = var.use_ipv6 == true ? 1 : 0
}

resource "aws_elb" "api_external" {
  name            = "${var.cluster_id}-ext"
  subnets         = data.aws_subnet.public.*.id
  internal        = false
  security_groups = ["${aws_security_group.master.id}", "${aws_security_group.api[0].id}"]

  idle_timeout                = 3600
  connection_draining         = true
  connection_draining_timeout = 300

  listener {
    instance_port     = 6443
    instance_protocol = "tcp"
    lb_port           = 6443
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "SSL:6443"
    interval            = 5
  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}-api-external",
  ), var.tags)}"

  count = local.public_endpoints && var.use_ipv6 == true ? 1 : 0
}

