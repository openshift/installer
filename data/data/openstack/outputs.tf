output "master_ips" {
  value = module.masters.master_ips
}

output "bootstrap_floating_ip" {
  value = module.bootstrap.bootstrap_floating_ip
}
