resource "aws_s3_bucket" "tectonic" {
  # This bucket name must match the CNAME
  # https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html#VirtualHostingCustomURLs
  bucket = "${var.tectonic_cluster_name}-ncg.${var.tectonic_base_domain}"

  acl = "private"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-tectonic",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${local.cluster_id}"
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
      "tectonicClusterID", "${local.cluster_id}"
    ), var.tectonic_aws_extra_tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}

data "archive_file" "assets" {
  type       = "zip"
  source_dir = "./generated/"

  # Because the archive_file provider is a data source, depends_on can't be
  # used to guarantee that the tectonic/bootkube modules have generated
  # all the assets on disk before trying to archive them. Instead, we use their
  # ID outputs, that are only computed once the assets have actually been
  # written to disk. We re-hash the IDs (or dedicated module outputs, like module.bootkube.content_hash)
  # to make the filename shorter, since there is no security nor collision risk anyways.
  #
  # Additionally, data sources do not support managing any lifecycle whatsoever,
  # and therefore, the archive is never deleted. To avoid cluttering the module
  # folder, we write it in the Terraform managed hidden folder `.terraform`.
  output_path = "./.terraform/generated_${sha1("${local.cluster_id}")}.zip"
}

# Ignition
resource "aws_s3_bucket_object" "ignition_bootstrap" {
  bucket  = "${aws_s3_bucket.tectonic.bucket}"
  key     = "ignition"
  content = "${file("./generated/ignition/bootstrap.json")}"
  acl     = "public-read"

  # TODO: Lock down permissions.
  # At the minute this is pulic (so accessible via http) so joiners nodes can reach the NCG using the same url
  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-ignition-master",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${local.cluster_id}"
    ), var.tectonic_aws_extra_tags)}"
}

resource "aws_s3_bucket_object" "ignition_etcd" {
  count   = "${length(data.template_file.etcd_hostname_list.*.id)}"
  bucket  = "${aws_s3_bucket.tectonic.bucket}"
  key     = "ignition_etcd_${count.index}.json"
  content = "${file("./generated/ignition/etcd-${count.index}.json")}"
  acl     = "private"

  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-ignition-etcd-${count.index}",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${local.cluster_id}"
    ), var.tectonic_aws_extra_tags)}"
}
