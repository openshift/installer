output "bootstrap_ip" {
  value = module.loadbalancer.powervs_lb_hostname
}

output "control_plane_ips" {
  value = module.master.master_ips
}
