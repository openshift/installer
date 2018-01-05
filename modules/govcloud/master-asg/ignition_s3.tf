resource "aws_s3_bucket_object" "ignition_master" {
  bucket  = "${var.s3_bucket}"
  key     = "ignition_master.json"
  content = "${data.ignition_config.main.rendered}"
  acl     = "public-read"

  server_side_encryption = "AES256"
}

data "ignition_config" "s3" {
  replace {
    source       = "${format("https://s3-us-gov-west-1.amazonaws.com/%s/%s", var.s3_bucket, aws_s3_bucket_object.ignition_master.key)}"
    verification = "sha512-${sha512(data.ignition_config.main.rendered)}"
  }
}
