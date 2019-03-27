output "master_ips" {
  value = "${split(" ", trimspace(data.external.ping.result.master_ips))}"
}

output "worker_ips" {
  value = "${split(" ", trimspace(data.external.ping.result.worker_ips))}"
}

output "bootstrap_ip" {
  value = "${split(" ", trimspace(data.external.ping.result.bootstrap_ip))}"
}
