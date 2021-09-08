output "bootstrap_a" {
  description = "IP Address of the bootstrap node"
  value       = data.dns_a_record_set.bootstrap.addrs[0]
}

output "lb_a" {
  description = "IP Address of the LoadBalancer node"
  value       = data.dns_a_record_set.lb.addrs[0]
}

output "masters_a" {
  description = "IP Addresses of the bootstrap node"
  // TODO: this assume 1 address per master
  value = flatten(data.dns_a_record_set.masters.*.addrs)
}

/*
output "workers_a" {
  description = "IP Addresses of the bootstrap node"
  value       = data.dns_a_record_set.workers.addrs
}
*/
