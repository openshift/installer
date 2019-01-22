output "ip_addresses" {
  value = "${google_compute_instance.master.*.network_interface.0.network_ip}"
}

output "master_instances" {
  value = "${google_compute_instance.master.*.self_link}"
}

output "master_instance_groups" {
  value = [
    "${google_compute_instance_group.master-0.self_link}",
    "${google_compute_instance_group.master-1.self_link}",
    "${google_compute_instance_group.master-2.self_link}",
  ]
}
