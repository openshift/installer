output "master_instances" {
  value = google_compute_instance.master.*.self_link
}

output "master_instance_groups" {
  value = google_compute_instance_group.master.*.self_link
}

output "ip_addresses" {
  value = google_compute_instance.master.*.network_interface.0.network_ip
}
