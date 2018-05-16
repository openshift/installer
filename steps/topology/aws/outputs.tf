# Etcd
output "etcd_sg_id" {
  value = "${module.vpc.etcd_sg_id}"
}

output "s3_bucket" {
  value = "${aws_s3_bucket.tectonic.bucket}"
}

# Masters
output "subnet_ids_masters" {
  value = "${module.vpc.master_subnet_ids}"
}

output "aws_lbs" {
  value = "${module.vpc.aws_lbs}"
}

output "master_sg_id" {
  value = "${module.vpc.master_sg_id}"
}

# Workers
output "subnet_ids_workers" {
  value = "${module.vpc.worker_subnet_ids}"
}

output "worker_sg_id" {
  value = "${module.vpc.worker_sg_id}"
}

# TNC
output "private_zone_id" {
  value = "${join("", aws_route53_zone.tectonic_int.*.zone_id)}"
}

output "tnc_elb_dns_name" {
  value = "${module.vpc.aws_elb_tnc_dns_name}"
}

output "tnc_elb_zone_id" {
  value = "${module.vpc.aws_elb_tnc_zone_id}"
}

output "tnc_s3_bucket_domain_name" {
  value = "${aws_s3_bucket.tectonic.bucket_domain_name}"
}
