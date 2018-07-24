data "ignition_config" "tnc" {
  count = "${var.instance_count}"

  append {
    source = "${format("http://${var.cluster_name}-tnc.${var.base_domain}/config/etcd?etcd_index=%d", count.index)}"

    # TODO: add verification
  }

  # Used for loading certificates
  append {
    source = "${format("s3://%s/ignition_etcd_%d.json", var.s3_bucket, count.index)}"

    # TODO: add verification
  }
}
