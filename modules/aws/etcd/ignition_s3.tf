locals {
  ignition_etcd_keys = ["ignition_etcd_0.json", "ignition_etcd_1.json", "ignition_etcd_2.json"]
}

data "ignition_config" "s3" {
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"

  replace {
    source = "${format("s3://%s/%s", var.s3_bucket, local.ignition_etcd_keys[count.index])}"

    # TODO: add verification
  }
}
