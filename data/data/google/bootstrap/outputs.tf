output "bootstrap_instance" {
  value = "${google_compute_instance.bootstrap.self_link}"
}

output "bootstrap_instance_group" {
  value = "${google_compute_instance_group.bootstrap.self_link}"
}
