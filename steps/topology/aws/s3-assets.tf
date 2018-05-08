## Create the S3 bucket where we'll upload the initial ignition file
# This needs a bit of a special setup. Because Ignition is fetching the 
# configuration over the tnc LB pointed to the s3 bucket, it doesn't have any
# identity. So, the file needs to be public. But then we'd expose secrets,
# so the public ignition file just has an ignition redirect to a s3:// url,
# which ignition can fetch directly with authentication.

resource "aws_s3_bucket" "tectonic" {
  # This bucket name must match the CNAME
  # https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html#VirtualHostingCustomURLs
  bucket = "${lower(var.tectonic_cluster_name)}-tnc.${var.tectonic_base_domain}"

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

# The real ignition contents, with secrets.
# Must be private. The node will zero this out as soon as it boots.
resource "aws_s3_bucket_object" "ignition_bootstrap_real" {
  bucket  = "${aws_s3_bucket.tectonic.bucket}"
  key     = "config/bootstrap"
  content = "${local.ignition_bootstrap}"
  acl     = "private"

  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-ignition-master",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}

# The public ignition configuration
data "ignition_config" "bootstrap_redirect" {
  replace {
    source = "s3://${aws_s3_bucket.tectonic.bucket}/config/bootstrap"
  }
}

# The public ignition object.
resource "aws_s3_bucket_object" "ignition_bootstrap" {
  bucket  = "${aws_s3_bucket.tectonic.bucket}"
  key     = "config/master"
  content = "${data.ignition_config.bootstrap_redirect.rendered}"
  acl     = "public-read"

  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-ignition-master",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}
