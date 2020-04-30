output "master_ips" {
  value = local.master_public_ipv4
}

output "worker_ips" {
  value = local.worker_public_ipv4
}

output "bootstrap_ip" {
  value = module.bootstrap.device_ip
}
