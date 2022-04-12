# Using explicit depends_on as otherwise there are issues with updating and adding of pool members
# Ref: https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/is_lb_listener

# even with the depends_on in machine_config_member, this still fails because the LB
# is in UPDATE_PENDING, so put in a time_sleep
resource "time_sleep" "wait_for_bootstrap_int_lb" {
  create_duration = "10s"
  depends_on      = [var.bootstrap_api_member_int_id]
}

resource "ibm_is_lb_pool_member" "machine_config_member" {
  count          = var.instance_count
  depends_on     = [time_sleep.wait_for_bootstrap_int_lb]
  lb             = var.lb_int_id
  pool           = var.machine_cfg_pool_id
  port           = 22623
  target_address = var.master_ips[count.index]
}

resource "ibm_is_lb_pool_member" "api_member_int" {
  count          = var.instance_count
  depends_on     = [ibm_is_lb_pool_member.machine_config_member]
  lb             = var.lb_int_id
  pool           = var.api_pool_int_id
  port           = 6443
  target_address = var.master_ips[count.index]
}

# even with the depends_on in api_member, this still fails because the LB
# is in UPDATE_PENDING, so put in a time_sleep
resource "time_sleep" "wait_for_bootstrap_ext_lb" {
  create_duration = "10s"
  depends_on      = [var.bootstrap_api_member_ext_id]
}

resource "ibm_is_lb_pool_member" "api_member" {
  count          = var.instance_count
  depends_on     = [time_sleep.wait_for_bootstrap_ext_lb]
  lb             = var.lb_ext_id
  pool           = var.api_pool_ext_id
  port           = 6443
  target_address = var.master_ips[count.index]
}
