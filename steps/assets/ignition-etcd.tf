locals {
  etcd_internal_instance_count = "${length(data.template_file.etcd_hostname_list.*.id)}"
  etcd_instance_count          = "${length(compact(var.tectonic_etcd_servers)) == 0 ? local.etcd_internal_instance_count : 0}"
}

data "ignition_config" "etcd" {
  count = "${local.etcd_instance_count}"

  files = ["${module.ignition_bootstrap.etcd_crt_id_list}"]
}
