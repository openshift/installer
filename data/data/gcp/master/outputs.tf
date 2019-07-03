output "master_instances" {
  value = google_compute_instance.master.*.self_link
}
