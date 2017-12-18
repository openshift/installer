data "aws_region" "current" {
  current = true
}

locals {
  govcloud_partition = "${data.aws_region.current.name == "us-gov-west-1" ? 1 : 0}"
  aws_partition      = "${data.aws_region.current.name == "us-gov-west-1" ? 0 : 1}"
  ami_owner          = "${data.aws_region.current.name == "us-gov-west-1" ?  "190570271432" : "595879546273"}"
  arn                = "${data.aws_region.current.name == "us-gov-west-1" ? "aws-us-gov" : "aws"}"
  s3_object_key      = "ignition_master.json"

  // Ignition does not support s3 protocol for GovCloud https://github.com/coreos/ignition/pull/477
  s3_endpoint = "${data.aws_region.current.name == "us-gov-west-1" ?
    format("https://s3-us-gov-west-1.amazonaws.com/%s/%s", var.s3_bucket, local.s3_object_key) :
    format("s3://%s/%s", var.s3_bucket, local.s3_object_key)}"
}

// Ignition config
resource "aws_s3_bucket_object" "ignition_master_govcloud" {
  count   = "${local.govcloud_partition}"
  bucket  = "${var.s3_bucket}"
  key     = "${local.s3_object_key}"
  content = "${data.ignition_config.main.rendered}"
  acl     = "public-read"

  server_side_encryption = "AES256"

  // Terraform does not support tags here for GovCloud
  // https://github.com/terraform-providers/terraform-provider-aws/pull/2665
}
