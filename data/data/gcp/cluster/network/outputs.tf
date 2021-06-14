output "cluster_ip" {
  value = google_compute_address.cluster_ip.address
}

output "cluster_public_ip" {
  value = var.public_endpoints ? google_compute_address.cluster_public_ip[0].address : null
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

output "api_health_checks" {
  value = google_compute_http_health_check.api.*.self_link
}

output "api_internal_health_checks" {
  value = google_compute_health_check.api_internal.*.self_link
}
