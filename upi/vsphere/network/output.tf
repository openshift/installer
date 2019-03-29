output "control_plane_ips" {
  value = "${split(" ", trimspace(data.external.ping.result.control_plane_ips))}"
}

output "compute_ips" {
  value = "${split(" ", trimspace(data.external.ping.result.compute_ips))}"
}

output "bootstrap_ip" {
  value = "${trimspace(data.external.ping.result.bootstrap_ip)}"
}
