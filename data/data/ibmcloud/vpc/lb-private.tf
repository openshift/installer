locals {
  # NOTE: Defined in ./vpc.tf
  # prefix = var.cluster_id

  # NOTE: Defined in ./lb-public.tf
  # port_ingress_http   = 80
  # port_ingress_https  = 443
  # port_kubernetes_api = 6443

  port_machine_config = 22623
}

############################################
# Load balancers
############################################

resource "ibm_is_lb" "kubernetes_api_private" {
  name            = "${local.prefix}-kubernetes-api-private"
  resource_group  = var.resource_group_id
  security_groups = [ ibm_is_security_group.control_plane.id ]
  subnets         = ibm_is_subnet.control_plane.*.id
  type            = "private"
}

############################################
# Load balancer backend pools
############################################

resource "ibm_is_lb_pool" "kubernetes_api_private" {
  name                = "${local.prefix}-kubernetes-api-private"
  lb                  = ibm_is_lb.kubernetes_api_private.id
  algorithm           = "round_robin"
  protocol            = "tcp"
  health_delay        = 60
  health_retries      = 5
  health_timeout      = 30
  health_type         = "https"
  health_monitor_url  = "/readyz"
  health_monitor_port = local.port_kubernetes_api
}

resource "ibm_is_lb_pool" "machine_config" {
  name                = "${local.prefix}-machine-config"
  lb                  = ibm_is_lb.kubernetes_api_private.id
  algorithm           = "round_robin"
  protocol            = "tcp"
  health_delay        = 60
  health_retries      = 5
  health_timeout      = 30
  health_type         = "tcp"
  health_monitor_port = local.port_machine_config
}

############################################
# Load balancer frontend listeners
############################################

resource "ibm_is_lb_listener" "kubernetes_api_private" {
  lb           = ibm_is_lb.kubernetes_api_private.id
  default_pool = ibm_is_lb_pool.kubernetes_api_private.id
  port         = local.port_kubernetes_api
  protocol     = "tcp"
}

resource "ibm_is_lb_listener" "machine_config" {
  lb           = ibm_is_lb.kubernetes_api_private.id
  default_pool = ibm_is_lb_pool.machine_config.id
  port         = local.port_machine_config
  protocol     = "tcp"
}
