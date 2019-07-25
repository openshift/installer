output "bootstrap_instances" {
  value = google_compute_instance.bootstrap.*.self_link
}

output "bootstrap_instance_ips" {
  value = google_compute_instance.bootstrap.*.network_interface.0.network_ip
}
