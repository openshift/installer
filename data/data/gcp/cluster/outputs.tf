output "cluster_ip" {
  value = module.network.cluster_ip
}

output "cluster_public_ip" {
  value = module.network.cluster_public_ip
}

output "network" {
  value = module.network.network
}

output "master_subnet" {
  value = module.network.master_subnet
}

output "api_health_checks" {
  value = module.network.api_health_checks
}

output "api_internal_health_checks" {
  value = module.network.api_internal_health_checks
}

output "master_instances" {
  value = module.master.master_instances
}

output "master_instance_groups" {
  value = module.master.master_instance_groups
}

output "compute_image" {
  value = local.gcp_image
}

output "control_plane_ips" {
  value = module.master.ip_addresses
}
