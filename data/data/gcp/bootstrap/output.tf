output "bootstrap_instances" {
  value = google_compute_instance.bootstrap.*.self_link
}
