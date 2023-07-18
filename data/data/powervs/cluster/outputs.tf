output "control_plane_ips" {
  value = module.master.master_ips
}

output "cluster_key_name" {
  value = ibm_pi_key.cluster_key.name
}

output "boot_image_id" {
  value = ibm_pi_image.boot_image.image_id
}

output "lb_int_id" {
  value = module.loadbalancer.lb_int_id
}

output "lb_ext_id" {
  value = module.loadbalancer.lb_ext_id
}

output "machine_cfg_pool_id" {
  value = module.loadbalancer.machine_cfg_pool_id
}

output "api_pool_int_id" {
  value = module.loadbalancer.api_pool_int_id
}

output "api_pool_ext_id" {
  value = module.loadbalancer.api_pool_ext_id
}

output "dhcp_id" {
  value = module.pi_network.dhcp_id
}

output "dhcp_network_id" {
  value = module.pi_network.dhcp_network_id
}

output "proxy_server_ip" {
  value = module.dns.dns_server
}

output "vpc_id" {
  value = module.vpc.vpc_id
}
