output "bootstrap_instances" {
  value = google_compute_instance.bootstrap.*.self_link
}

output "bootstrap_instance_groups" {
  value = google_compute_instance_group.bootstrap.*.self_link
}

output "ip_addresses" {
  value = google_compute_instance.bootstrap[0].network_interface.0.network_ip
}

