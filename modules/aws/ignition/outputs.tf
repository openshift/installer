output "ignition" {
  value = "${data.ignition_config.main.rendered}"
}
