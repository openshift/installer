locals {
  # NOTE: Defined in ./vpc.tf
  # prefix = var.cluster_id

  port_kubernetes_api = 6443
}

############################################
# Load balancers
############################################

resource "ibm_is_lb" "kubernetes_api_public" {
  count           = var.public_endpoints ? 1 : 0

  name            = "${local.prefix}-kubernetes-api-public"
  resource_group  = var.resource_group_id
  security_groups = [ ibm_is_security_group.control_plane.id ]
  subnets         = ibm_is_subnet.control_plane.*.id
  tags            = var.tags
  type            = "public"
}

############################################
# Load balancer backend pools
############################################

resource "ibm_is_lb_pool" "kubernetes_api_public" {
  count               = var.public_endpoints ? 1 : 0

  name                = "${local.prefix}-kubernetes-api-public"
  lb                  = ibm_is_lb.kubernetes_api_public.0.id
  algorithm           = "round_robin"
  protocol            = "tcp"
  health_delay        = 60
  health_retries      = 5
  health_timeout      = 30
  health_type         = "https"
  health_monitor_url  = "/readyz"
  health_monitor_port = local.port_kubernetes_api
}

############################################
# Load balancer frontend listeners
############################################

resource "ibm_is_lb_listener" "kubernetes_api_public" {
  count        = var.public_endpoints ? 1 : 0

  lb           = ibm_is_lb.kubernetes_api_public.0.id
  default_pool = ibm_is_lb_pool.kubernetes_api_public.0.id
  port         = local.port_kubernetes_api
  protocol     = "tcp"
}
