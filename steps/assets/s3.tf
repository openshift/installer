data "aws_region" "current" {
  current = true
}

resource "aws_s3_bucket" "tectonic" {
  # This name must match the CNAME
  # https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html#VirtualHostingCustomURLs
  bucket = "${var.tectonic_cluster_name}-ncg.${var.tectonic_base_domain}"

  acl = "private"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-tectonic",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${module.tectonic.cluster_id}"
    ), var.tectonic_aws_extra_tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}

# Bootkube / Tectonic assets
resource "aws_s3_bucket_object" "tectonic_assets" {
  bucket = "${aws_s3_bucket.tectonic.bucket}"
  key    = "assets.zip"
  source = "${data.archive_file.assets.output_path}"
  acl    = "private"

  # To be on par with the current Tectonic installer, we only do server-side
  # encryption, using AES256. Eventually, we should start using KMS-based
  # client-side encryption.
  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-tectonic-assets",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${module.tectonic.cluster_id}"
    ), var.tectonic_aws_extra_tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}

# kubeconfig
resource "aws_s3_bucket_object" "kubeconfig" {
  bucket  = "${aws_s3_bucket.tectonic.bucket}"
  key     = "kubeconfig"
  content = "${module.bootkube.kubeconfig-kubelet}"
  acl     = "private"

  # The current Tectonic installer stores bits of the kubeconfig in KMS. As we
  # do not support KMS yet, we at least offload it to S3 for now. Eventually,
  # we should consider using KMS-based client-side encryption, or uploading it
  # to KMS.
  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-kubeconfig",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${module.tectonic.cluster_id}"
    ), var.tectonic_aws_extra_tags)}"
}
