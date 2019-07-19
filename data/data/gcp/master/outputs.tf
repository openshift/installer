output "master_instances" {
  value = google_compute_instance.master.*.self_link
}

output "ip_addresses" {
  value = google_compute_instance.master.*.network_interface.0.network_ip
}
