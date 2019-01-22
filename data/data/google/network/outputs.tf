output "google_lb_api_external_address" {
  value = "${google_compute_global_address.master.address}"
}

output "network" {
  value = "${google_compute_network.default.self_link}"
}

output "subnetwork" {
  value = "${google_compute_subnetwork.default.self_link}"
}

output "zones" {
  value = "${local.zones}"
}
