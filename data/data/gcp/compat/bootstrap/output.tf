output "bootstrap_instances" {
  value = google_compute_instance.bootstrap.*.self_link
}

output "bootstrap_instance_groups" {
  value = google_compute_instance_group.bootstrap.*.self_link
}
