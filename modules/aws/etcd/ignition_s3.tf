resource "aws_s3_bucket_object" "ignition_etcd" {
  count = "${local.aws_partition == 1 ? (length(var.external_endpoints) == 0 ? var.instance_count : 0) : 0}"

  bucket  = "${var.s3_bucket}"
  key     = "ignition_etcd_${count.index}.json"
  content = "${data.ignition_config.etcd.*.rendered[count.index]}"
  acl     = "private"

  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.cluster_name}-ignition-etcd-${count.index}",
      "KubernetesCluster", "${var.cluster_name}",
      "tectonicClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

data "ignition_config" "s3" {
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"

  replace {
    source       = "${local.s3_endpoints[count.index]}"
    verification = "sha512-${sha512(data.ignition_config.etcd.*.rendered[count.index])}"
  }
}
