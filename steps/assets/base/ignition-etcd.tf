data "ignition_config" "etcd" {
  count = "${var.etcd_count}"

  files = ["${module.ignition_bootstrap.etcd_crt_id_list}"]
}
