data "aws_lb" "api_internal" {
  name = "${var.cluster_id}-int"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-int"
    },
    var.tags,
  )
}

data "aws_lb" "api_external" {
  count = local.public_endpoints ? 1 : 0
  name  = "${var.cluster_id}-ext"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-ext"
    },
    var.tags,
  )
}

data "aws_lb_target_group" "api_internal" {
  name = "${var.cluster_id}-aint"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-aint"
    },
    var.tags,
  )
}

data "aws_lb_target_group" "api_external" {
  count = local.public_endpoints ? 1 : 0
  name  = "${var.cluster_id}-aext"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-aext"
    },
    var.tags,
  )
}

data "aws_lb_target_group" "services" {
  name = "${var.cluster_id}-sint"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-sint"
    },
    var.tags,
  )
}

data "aws_lb_listener" "api_internal_api" {
  load_balancer_arn = data.aws_lb.api_internal.arn
  port              = "6443"
}

data "aws_lb_listener" "api_internal_services" {
  load_balancer_arn = data.aws_lb.api_internal.arn
  port              = "22623"
}

data "aws_lb_listener" "api_external_api" {
  count = local.public_endpoints ? 1 : 0

  load_balancer_arn = data.aws_lb.api_external[0].arn
  port              = "6443"
}

