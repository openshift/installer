output "max_user_watches_id" {
  value = "${data.ignition_file.max_user_watches.id}"
}

output "max_user_watches_rendered" {
  value = "${data.template_file.max_user_watches.rendered}"
}

output "docker_dropin_id" {
  value = "${data.ignition_systemd_unit.docker_dropin.id}"
}

output "docker_dropin_rendered" {
  value = "${data.template_file.docker_dropin.rendered}"
}

output "kubelet_service_id" {
  value = "${data.ignition_systemd_unit.kubelet.id}"
}

output "kubelet_service_rendered" {
  value = "${data.template_file.kubelet.rendered}"
}

output "kubelet_env_service_id" {
  value = "${data.ignition_systemd_unit.kubelet_env.id}"
}

output "kubelet_env_service_rendered" {
  value = "${data.template_file.kubelet_env_service.rendered}"
}

output "s3_puller_id" {
  value = "${data.ignition_file.s3_puller.id}"
}

output "s3_puller_rendered" {
  value = "${data.template_file.s3_puller.rendered}"
}

output "locksmithd_service_id" {
  value = "${data.ignition_systemd_unit.locksmithd.id}"
}

output "kubelet_env_id" {
  value = "${data.ignition_file.kubelet_env.id}"
}

output "kubelet_env_rendered" {
  value = "${data.template_file.kubelet_env.rendered}"
}

output "tx_off_service_id" {
  value = "${data.ignition_systemd_unit.tx_off.id}"
}

output "tx_off_service_rendered" {
  value = "${data.template_file.tx_off.rendered}"
}

output "azure_udev_rules_id" {
  value = "${data.ignition_file.azure_udev_rules.id}"
}

output "azure_udev_rules_rendered" {
  value = "${data.template_file.azure_udev_rules.rendered}"
}
