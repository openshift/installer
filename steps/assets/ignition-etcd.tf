locals {
  etcd_internal_instance_count = "${length(data.template_file.etcd_hostname_list.*.id)}"
  etcd_instance_count          = "${length(compact(var.tectonic_etcd_servers)) == 0 ? local.etcd_internal_instance_count : 0}"
}

resource "aws_s3_bucket_object" "ignition_etcd" {
  count   = "${local.etcd_instance_count}"
  bucket  = "${aws_s3_bucket.tectonic.bucket}"
  key     = "ignition_etcd_${count.index}.json"
  content = "${data.ignition_config.etcd.*.rendered[count.index]}"
  acl     = "private"

  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-ignition-etcd-${count.index}",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${module.tectonic.cluster_id}"
    ), var.tectonic_aws_extra_tags)}"
}

data "ignition_config" "etcd" {
  count = "${local.etcd_instance_count}"

  systemd = [
    "${data.ignition_systemd_unit.locksmithd.*.id[count.index]}",
    "${module.ignition_bootstrap.etcd_dropin_id_list[count.index]}",
  ]

  files = ["${compact(list(
    module.ignition_bootstrap.profile_env_id,
    module.ignition_bootstrap.systemd_default_env_id,
   ))}",
    "${module.ignition_bootstrap.etcd_crt_id_list}",
  ]
}

data "ignition_systemd_unit" "locksmithd" {
  count = "${local.etcd_instance_count}"

  name    = "locksmithd.service"
  enabled = true

  dropin = [
    {
      name = "40-etcd-lock.conf"

      content = <<EOF
[Service]
Environment=REBOOT_STRATEGY=etcd-lock
Environment="LOCKSMITHD_ETCD_CAFILE=/etc/ssl/etcd/ca.crt"
Environment="LOCKSMITHD_ETCD_KEYFILE=/etc/ssl/etcd/client.key"
Environment="LOCKSMITHD_ETCD_CERTFILE=/etc/ssl/etcd/client.crt"
Environment="LOCKSMITHD_ENDPOINT=https://${var.tectonic_cluster_name}-etcd-${count.index}.${var.tectonic_base_domain}:2380"
EOF
    },
  ]
}
