output "lb_hostname" {
  value = ibm_is_lb.load_balancer.hostname
}

output "lb_int_hostname" {
  value = ibm_is_lb.load_balancer_int.hostname
}

output "lb_ext_id" {
  value = ibm_is_lb.load_balancer.id
}

output "machine_cfg_pool_id" {
  value = ibm_is_lb_pool.machine_config_pool.id
}

output "lb_int_id" {
  value = ibm_is_lb.load_balancer_int.id
}

output "api_pool_int_id" {
  value = ibm_is_lb_pool.api_pool_int.id
}

output "api_pool_ext_id" {
  value = ibm_is_lb_pool.api_pool.id
}
