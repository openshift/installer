resource "aws_s3_bucket_object" "ignition_etcd" {
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"

  bucket  = "${var.s3_bucket}"
  key     = "ignition_etcd_${count.index}.json"
  content = "${data.ignition_config.etcd.*.rendered[count.index]}"
  acl     = "public-read"

  server_side_encryption = "AES256"
}

data "ignition_config" "s3" {
  count = "${length(var.external_endpoints) == 0 ? var.instance_count : 0}"

  replace {
    source       = "${format("https://s3-us-gov-west-1.amazonaws.com/%s/%s", var.s3_bucket, aws_s3_bucket_object.ignition_etcd.*.key[count.index])}"
    verification = "sha512-${sha512(data.ignition_config.etcd.*.rendered[count.index])}"
  }
}
