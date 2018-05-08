output "ignition_bootstrap" {
  value = "${data.ignition_config.bootstrap.rendered}"
}

output "ignition_etcd" {
  value = "${module.assets_base.ignition_etcd}"
}
