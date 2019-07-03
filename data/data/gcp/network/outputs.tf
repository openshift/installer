output "network" {
  value = google_compute_network.cluster_network.self_link
}

output "worker_subnet" {
  value = google_compute_subnetwork.worker_subnet.self_link
}

output "master_subnet" {
  value = google_compute_subnetwork.master_subnet.self_link
}

output "zones" {
  value = local.zones
}
