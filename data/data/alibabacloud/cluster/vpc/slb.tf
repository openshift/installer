
resource "alicloud_slb_load_balancer" "slb_external" {
  count = local.is_external ? 1 : 0

  resource_group_id  = var.resource_group_id
  load_balancer_name = "${local.prefix}-slb-external"
  address_type       = "internet"
  load_balancer_spec = "slb.s2.small"

  tags = merge(
    {
      "Name" = "${local.prefix}-slb-external"
    },
    var.tags,
  )
}

resource "alicloud_slb_listener" "listener_external_6443" {
  count = local.is_external ? 1 : 0

  load_balancer_id          = alicloud_slb_load_balancer.slb_external[0].id
  backend_port              = 6443
  frontend_port             = 6443
  protocol                  = "tcp"
  bandwidth                 = -1
  sticky_session            = "off"
  cookie_timeout            = 86400
  health_check              = "on"
  health_check_type         = "http"
  health_check_uri          = "/readyz"
  health_check_connect_port = 6080
  healthy_threshold         = 2
  unhealthy_threshold       = 2
  health_check_timeout      = 10
  health_check_interval     = 10
  x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  request_timeout = 80
  idle_timeout    = 30
}

resource "alicloud_slb_load_balancer" "slb_internal" {
  resource_group_id  = var.resource_group_id
  load_balancer_name = "${local.prefix}-slb-internal"
  address_type       = "intranet"
  vswitch_id         = local.vswitch_ids[0]
  load_balancer_spec = "slb.s2.small"
  tags = merge(
    {
      "Name" = "${local.prefix}-slb-internal"
    },
    var.tags,
  )
}

resource "alicloud_slb_listener" "listener_internal_6443" {
  load_balancer_id          = alicloud_slb_load_balancer.slb_internal.id
  backend_port              = 6443
  frontend_port             = 6443
  protocol                  = "tcp"
  bandwidth                 = -1
  sticky_session            = "off"
  cookie_timeout            = 86400
  health_check              = "on"
  health_check_type         = "http"
  health_check_uri          = "/readyz"
  health_check_connect_port = 6080
  healthy_threshold         = 2
  unhealthy_threshold       = 2
  health_check_timeout      = 10
  health_check_interval     = 10
  x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  request_timeout = 80
  idle_timeout    = 30
}

resource "alicloud_slb_listener" "listener_internal_22623" {
  load_balancer_id          = alicloud_slb_load_balancer.slb_internal.id
  backend_port              = 22623
  frontend_port             = 22623
  protocol                  = "tcp"
  bandwidth                 = -1
  sticky_session            = "off"
  cookie_timeout            = 86400
  health_check              = "on"
  health_check_type         = "http"
  health_check_uri          = "/healthz"
  health_check_connect_port = 22624
  healthy_threshold         = 2
  unhealthy_threshold       = 2
  health_check_timeout      = 10
  health_check_interval     = 10
  x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  request_timeout = 80
  idle_timeout    = 30
}
