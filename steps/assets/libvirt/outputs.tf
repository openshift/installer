output "ignition_bootstrap" {
  value = "${data.ignition_config.bootstrap.rendered}"
}

output "ignition_etcd" {
  value = "${data.ignition_config.etcd.*.rendered}"
}
