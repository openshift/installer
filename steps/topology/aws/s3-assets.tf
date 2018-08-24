resource "aws_s3_bucket" "tectonic" {
  bucket = "${lower(var.tectonic_cluster_name)}.${var.tectonic_base_domain}"

  acl = "private"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-tectonic",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}
