output "master_ips" {
  value = ["${local.master_public_ipv4}"]
}

output "worker_ips" {
  value = ["${local.master_public_ipv4}"]
}

output "bootstrap_ip" {
  value = "${module.bootstrap.device_ip}"
}

output "sos_consoles" {
  value = "${zipmap(concat(list(module.bootstrap.device_hostname), packet_device.masters.*.hostname, packet_device.workers.*.hostname), formatlist("ssh %s@sos.%s.packet.net", concat(list(module.bootstrap.device_id), packet_device.masters.*.id, packet_device.workers.*.id), local.packet_facility))}"
}
