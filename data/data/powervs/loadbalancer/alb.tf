locals {
  api_servers       = concat([var.bootstrap_ip], var.master_ips)
  api_servers_count = length(var.master_ips) + 1 # bootstrap + master
  app_servers       = var.master_ips
  app_servers_count = length(var.master_ips)
}

resource "ibm_is_lb" "load_balancer" {
  name            = "${var.cluster_id}-loadbalancer"
  subnets         = [var.vpc_subnet_id]
  security_groups = [ibm_is_security_group.ocp_security_group.id]
  tags            = [var.cluster_id]
}

# Using explicit depends_on as otherwise there are issues with updating and adding of pool members
# Ref: https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/is_lb_listener

## TODO move this to internal/private LB
# machine config listener and backend pool
resource "ibm_is_lb_listener" "machine_config_listener" {
  lb           = ibm_is_lb.load_balancer.id
  port         = 22623
  protocol     = "tcp"
  default_pool = ibm_is_lb_pool.machine_config_pool.id
}
resource "ibm_is_lb_pool" "machine_config_pool" {
  depends_on = [ibm_is_lb.load_balancer]

  name           = "machine-config-server"
  lb             = ibm_is_lb.load_balancer.id
  algorithm      = "round_robin"
  protocol       = "tcp"
  health_delay   = 60
  health_retries = 5
  health_timeout = 30
  health_type    = "tcp"
}
resource "ibm_is_lb_pool_member" "machine_config_member" {
  depends_on = [ibm_is_lb_listener.machine_config_listener]
  count      = local.api_servers_count

  lb             = ibm_is_lb.load_balancer.id
  pool           = ibm_is_lb_pool.machine_config_pool.id
  port           = 22623
  target_address = local.api_servers[count.index]
}


# api listener and backend pool
resource "ibm_is_lb_listener" "api_listener" {
  lb           = ibm_is_lb.load_balancer.id
  port         = 6443
  protocol     = "tcp"
  default_pool = ibm_is_lb_pool.api_pool.id
}
resource "ibm_is_lb_pool" "api_pool" {
  depends_on = [ibm_is_lb.load_balancer]

  name           = "openshift-api-server"
  lb             = ibm_is_lb.load_balancer.id
  algorithm      = "round_robin"
  protocol       = "tcp"
  health_delay   = 60
  health_retries = 5
  health_timeout = 30
  health_type    = "tcp"
}
resource "ibm_is_lb_pool_member" "api_member" {
  depends_on = [ibm_is_lb_listener.api_listener, ibm_is_lb_pool_member.machine_config_member]
  count      = local.api_servers_count

  lb             = ibm_is_lb.load_balancer.id
  pool           = ibm_is_lb_pool.api_pool.id
  port           = 6443
  target_address = local.api_servers[count.index]
}

