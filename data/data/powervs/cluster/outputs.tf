output "bootstrap_ip" {
  value = module.bootstrap.bootstrap_ip
}

output "control_plane_ips" {
  value = module.master.master_ips
}
