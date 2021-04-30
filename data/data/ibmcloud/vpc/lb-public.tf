locals {
  # NOTE: Defined in ./vpc.tf
  # prefix = var.cluster_id

  port_ingress_http   = 80
  port_ingress_https  = 443
  port_kubernetes_api = 6443
}

############################################
# Load balancers
############################################

resource "ibm_is_lb" "kubernetes_api_public" {
  name            = "${local.prefix}-kubernetes-api-public"
  resource_group  = var.resource_group_id
  security_groups = [ ibm_is_security_group.control_plane.id ]
  subnets         = ibm_is_subnet.control_plane.*.id
  type            = "public"
}

resource "ibm_is_lb" "application_ingress_public" {
  name            = "${local.prefix}-application-ingress-public"
  resource_group  = var.resource_group_id
  security_groups = [ ibm_is_security_group.compute.id ]
  subnets         = ibm_is_subnet.compute.*.id
  type            = "public"
}

############################################
# Load balancer backend pools
############################################

resource "ibm_is_lb_pool" "kubernetes_api_public" {
  name                = "${local.prefix}-kubernetes-api-public"
  lb                  = ibm_is_lb.kubernetes_api_public.id
  algorithm           = "round_robin"
  protocol            = "tcp"
  health_delay        = 60
  health_retries      = 5
  health_timeout      = 30
  health_type         = "https"
  health_monitor_url  = "/readyz"
  health_monitor_port = local.port_kubernetes_api
}

resource "ibm_is_lb_pool" "application_ingress_public_http" {
  name                = "${local.prefix}-application-ingress-public-http"
  lb                  = ibm_is_lb.application_ingress_public.id
  algorithm           = "round_robin"
  protocol            = "tcp"
  health_delay        = 60
  health_retries      = 5
  health_timeout      = 30
  health_type         = "tcp"
  health_monitor_port = local.port_ingress_http
}

resource "ibm_is_lb_pool" "application_ingress_public_https" {
  name                = "${local.prefix}-application-ingress-public-https"
  lb                  = ibm_is_lb.application_ingress_public.id
  algorithm           = "round_robin"
  protocol            = "tcp"
  health_delay        = 60
  health_retries      = 5
  health_timeout      = 30
  health_type         = "tcp"
  health_monitor_port = local.port_ingress_https
}

############################################
# Load balancer frontend listeners
############################################

resource "ibm_is_lb_listener" "kubernetes_api_public" {
  lb           = ibm_is_lb.kubernetes_api_public.id
  default_pool = ibm_is_lb_pool.kubernetes_api_public.id
  port         = local.port_kubernetes_api
  protocol     = "tcp"
}

resource "ibm_is_lb_listener" "application_ingress_public_http" {
  lb           = ibm_is_lb.application_ingress_public.id
  default_pool = ibm_is_lb_pool.application_ingress_public_http.id
  port         = local.port_ingress_http
  protocol     = "tcp"
}

resource "ibm_is_lb_listener" "application_ingress_public_https" {
  lb           = ibm_is_lb.application_ingress_public.id
  default_pool = ibm_is_lb_pool.application_ingress_public_https.id
  port         = local.port_ingress_https
  protocol     = "tcp"
}