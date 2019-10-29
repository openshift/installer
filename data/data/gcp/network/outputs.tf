output "cluster_ip" {
  value = google_compute_forwarding_rule.api_internal.ip_address
}

output "cluster_public_ip" {
  value = var.public_endpoints ? google_compute_forwarding_rule.api[0].ip_address : null
}

output "network" {
  value = local.cluster_network
}

output "worker_subnet" {
  value = var.preexisting_network ? data.google_compute_subnetwork.preexisting_worker_subnet[0].self_link : google_compute_subnetwork.worker_subnet[0].self_link
}

output "master_subnet" {
  value = local.master_subnet
}
