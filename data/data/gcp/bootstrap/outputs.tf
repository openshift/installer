output "bootstrap_instances" {
  value = google_compute_instance.bootstrap.*.self_link
}

output "bootstrap_instance_groups" {
  value = google_compute_instance_group.bootstrap.*.self_link
}

output "bootstrap_ip" {
  value = local.public_endpoints ? google_compute_instance.bootstrap.network_interface.0.access_config.0.nat_ip : google_compute_instance.bootstrap.network_interface.0.network_ip
}
