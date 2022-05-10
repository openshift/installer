# Using explicit depends_on as otherwise there are issues with updating and adding of pool members
# Ref: https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/is_lb_listener

# bootstrap ssh pool, listener, member
resource "ibm_is_lb_pool" "bootstrap_pool" {
  name           = "bootstrap-node"
  lb             = var.lb_ext_id
  algorithm      = "round_robin"
  protocol       = "tcp"
  health_delay   = 5
  health_retries = 2
  health_timeout = 2
  health_type    = "tcp"
}

# explicit depends because the LB will be in `UPDATE_PENDING` state and this will fail
resource "ibm_is_lb_listener" "bootstrap_listener" {
  lb           = var.lb_ext_id
  port         = 22
  protocol     = "tcp"
  default_pool = ibm_is_lb_pool.bootstrap_pool.id
}

resource "ibm_is_lb_pool_member" "bootstrap" {
  depends_on     = [ibm_is_lb_listener.bootstrap_listener]
  lb             = var.lb_ext_id
  pool           = ibm_is_lb_pool.bootstrap_pool.id
  port           = 22
  target_address = var.bootstrap_ip
}
