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
}

resource "aws_lb" "api_external" {
  count = local.public_endpoints ? 1 : 0

  name                             = "${var.cluster_id}-ext"
  load_balancer_type               = "network"
  internal                         = false
  enable_cross_zone_load_balancing = true

  dynamic "subnet_mapping" {
    for_each = range(length(data.aws_subnet.public))

    content {
      subnet_id     = data.aws_subnet.public[subnet_mapping.key].id
      allocation_id = aws_eip.api_nlb_public[subnet_mapping.key].id
    }
  }

  tags = merge(
    {
      "Name" = "${var.cluster_id}-ext"
    },
    var.tags,
  )

  timeouts {
    create = "20m"
  }
}

resource "aws_eip" "api_nlb_public" {
  count  = length(var.availability_zones)
  domain = "vpc"

  public_ipv4_pool = var.public_ipv4_pool == "" ? null : var.public_ipv4_pool

  tags = merge(
    {
      "Name" = "${var.cluster_id}-eip-${var.availability_zones[count.index]}-lb-api"
    },
    var.tags,
  )

  # Terraform does not declare an explicit dependency towards the internet gateway.
  # this can cause the internet gateway to be deleted/detached before the EIPs.
  # https://github.com/coreos/tectonic-installer/issues/1017#issuecomment-307780549
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

  // Refer to docs/dev/kube-apiserver-health-check.md on how to correctly setup health check probe for kube-apiserver
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
  count = local.public_endpoints ? 1 : 0

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

  // Refer to docs/dev/kube-apiserver-health-check.md on how to correctly setup health check probe for kube-apiserver
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
  count = local.public_endpoints ? 1 : 0

  load_balancer_arn = aws_lb.api_external[0].arn
  protocol          = "TCP"
  port              = "6443"

  default_action {
    target_group_arn = aws_lb_target_group.api_external[0].arn
    type             = "forward"
  }
}

