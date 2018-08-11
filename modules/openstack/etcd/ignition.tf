data "ignition_config" "tnc" {
  count = "${var.instance_count}"

  append {
    source = "${format("http://${var.cluster_name}-tnc.${var.base_domain}/config/etcd?etcd_index=%d", count.index)}"

    # TODO: add verification
  }

  # TODO: Figure out how to handle certificates
}
