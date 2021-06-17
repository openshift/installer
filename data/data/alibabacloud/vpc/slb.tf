
resource "alicloud_slb_load_balancer" "slb_external" {
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

resource "alicloud_slb_listener" "listener_external_80" {
  load_balancer_id    = alicloud_slb_load_balancer.slb_external.id
  backend_port        = 80
  frontend_port       = 80
  protocol            = "tcp"
  bandwidth           = 10
  sticky_session      = "on"
  sticky_session_type = "insert"
  cookie_timeout      = 86400
  # health_check              = "on"
  # health_check_domain       = "ali.com"
  # health_check_uri          = "/cons"
  # health_check_connect_port = 20
  # healthy_threshold         = 8
  # unhealthy_threshold       = 8
  # health_check_timeout      = 8
  # health_check_interval     = 5
  # health_check_http_code    = "http_2xx,http_3xx"
  x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  # acl_status      = "on"
  # acl_type        = "white"
  # acl_id          = alicloud_slb_acl.default.id
  request_timeout = 80
  idle_timeout    = 30
}


resource "alicloud_slb_listener" "listener_external_443" {
  load_balancer_id    = alicloud_slb_load_balancer.slb_external.id
  backend_port        = 443
  frontend_port       = 443
  protocol            = "tcp"
  bandwidth           = 10
  sticky_session      = "on"
  sticky_session_type = "insert"
  cookie_timeout      = 86400
  # health_check              = "on"
  # health_check_domain       = "ali.com"
  # health_check_uri          = "/cons"
  # health_check_connect_port = 20
  # healthy_threshold         = 8
  # unhealthy_threshold       = 8
  # health_check_timeout      = 8
  # health_check_interval     = 5
  # health_check_http_code    = "http_2xx,http_3xx"
  x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  # acl_status      = "on"
  # acl_type        = "white"
  # acl_id          = alicloud_slb_acl.default.id
  request_timeout = 80
  idle_timeout    = 30
}

resource "alicloud_slb_listener" "listener_external_6443" {
  load_balancer_id    = alicloud_slb_load_balancer.slb_external.id
  backend_port        = 6443
  frontend_port       = 6443
  protocol            = "tcp"
  bandwidth           = 10
  sticky_session      = "on"
  sticky_session_type = "insert"
  cookie_timeout      = 86400
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
  vswitch_id         = alicloud_vswitch.vswitchs[0].id
  load_balancer_spec = "slb.s2.small"
  tags = merge(
    {
      "Name" = "${local.prefix}-slb-internal"
    },
    var.tags,
  )
}

resource "alicloud_slb_listener" "listener_internal_6443" {
  load_balancer_id    = alicloud_slb_load_balancer.slb_internal.id
  backend_port        = 6443
  frontend_port       = 6443
  protocol            = "tcp"
  bandwidth           = 10
  sticky_session      = "on"
  sticky_session_type = "insert"
  cookie_timeout      = 86400
  x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  request_timeout = 80
  idle_timeout    = 30
}

resource "alicloud_slb_listener" "listener_internal_22623" {
  load_balancer_id    = alicloud_slb_load_balancer.slb_internal.id
  backend_port        = 22623
  frontend_port       = 22623
  protocol            = "tcp"
  bandwidth           = 10
  sticky_session      = "on"
  sticky_session_type = "insert"
  cookie_timeout      = 86400
  x_forwarded_for {
    retrive_slb_ip = true
    retrive_slb_id = true
  }
  request_timeout = 80
  idle_timeout    = 30
}
