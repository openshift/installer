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

output "k8s_node_bootstrap_service_id" {
  value = "${data.ignition_systemd_unit.k8s_node_bootstrap.id}"
}

output "k8s_node_bootstrap_service_rendered" {
  value = "${data.template_file.k8s_node_bootstrap.rendered}"
}

output "init_assets_service_id" {
  value = "${data.ignition_systemd_unit.init_assets.id}"
}

output "s3_puller_id" {
  value = "${data.ignition_file.s3_puller.id}"
}

output "s3_puller_rendered" {
  value = "${data.template_file.s3_puller.rendered}"
}

output "gcs_puller_id" {
  value = "${data.ignition_file.gcs_puller.id}"
}

output "gcs_puller_rendered" {
  value = "${data.template_file.gcs_puller.rendered}"
}

output "locksmithd_service_id" {
  value = "${data.ignition_systemd_unit.locksmithd.id}"
}

output "installer_kubelet_env_id" {
  value = "${data.ignition_file.installer_kubelet_env.id}"
}

output "installer_kubelet_env_rendered" {
  value = "${data.template_file.installer_kubelet_env.rendered}"
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

output "etcd_dropin_id_list" {
  value = "${data.ignition_systemd_unit.etcd.*.id}"
}

output "etcd_dropin_rendered_list" {
  value = "${data.template_file.etcd.*.rendered}"
}

output "coreos_metadata_dropin_id" {
  value = "${data.ignition_systemd_unit.coreos_metadata.id}"
}

output "coreos_metadata_dropin_rendered" {
  value = "${data.template_file.coreos_metadata.rendered}"
}
