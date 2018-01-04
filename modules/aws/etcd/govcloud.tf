data "aws_region" "current" {
  current = true
}

locals {
  govcloud_partition = "${data.aws_region.current.name == "us-gov-west-1" ? 1 : 0}"
  aws_partition      = "${data.aws_region.current.name == "us-gov-west-1" ? 0 : 1}"
  ami_owner          = "${data.aws_region.current.name == "us-gov-west-1" ?  "190570271432" : "595879546273"}"
  arn                = "${data.aws_region.current.name == "us-gov-west-1" ? "aws-us-gov" : "aws"}"

  s3_object_keys = "${data.aws_region.current.name == "us-gov-west-1" ?
    join(" ", aws_s3_bucket_object.ignition_etcd_govcloud.*.key) :
    join(" ", aws_s3_bucket_object.ignition_etcd.*.key)}"

  // Ignition does not support s3 protocol for GovCloud https://github.com/coreos/ignition/pull/477
  s3_prefix    = "${data.aws_region.current.name == "us-gov-west-1" ? "https://s3-us-gov-west-1.amazonaws.com" : "s3:/" }"
  s3_endpoints = "${formatlist("${local.s3_prefix}/%s/%s", var.s3_bucket, split(" ", local.s3_object_keys))}"
}

// Ignition config
resource "aws_s3_bucket_object" "ignition_etcd_govcloud" {
  count = "${local.govcloud_partition == 1 ? (length(var.external_endpoints) == 0 ? var.instance_count : 0) : 0}"

  bucket  = "${var.s3_bucket}"
  key     = "ignition_etcd_${count.index}.json"
  content = "${data.ignition_config.etcd.*.rendered[count.index]}"
  acl     = "public-read"

  server_side_encryption = "AES256"

  // Terraform does not support tags here for GovCloud
  // https://github.com/terraform-providers/terraform-provider-aws/pull/2665
}
