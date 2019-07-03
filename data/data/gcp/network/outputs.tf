output "network" {
  value = google_compute_network.cluster_network.self_link
}

output "compute_subnet" {
  value = google_compute_subnetwork.compute_subnet.self_link
}

output "control_subnet" {
  value = google_compute_subnetwork.control_subnet.self_link
}

output "zones" {
  value = local.zones
}
