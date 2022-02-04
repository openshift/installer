output "powervs_lb_hostname" {
  value = ibm_is_lb.load_balancer.hostname
}

output "powervs_lb_int_hostname" {
  value = ibm_is_lb.load_balancer_int.hostname
}
