output "subnet_ids_masters" {
  value = "${module.vpc.master_subnet_ids}"
}

output "aws_lbs" {
  value = "${module.vpc.aws_lbs}"
}

output "master_sg_id" {
  value = "${module.vpc.master_sg_id}"
}

output "s3_bucket" {
  value = "${aws_s3_bucket.tectonic.bucket}"
}
