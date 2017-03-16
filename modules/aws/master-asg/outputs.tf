output "api-internal-elb" {
  value = {
    dns_name = "${aws_elb.api-internal.dns_name}"
    zone_id  = "${aws_elb.api-internal.zone_id}"
  }
}

output "api-external-elb" {
  value = {
    dns_name = "${aws_elb.api-external.dns_name}"
    zone_id  = "${aws_elb.api-external.zone_id}"
  }
}

output "console-elb" {
  value = {
    dns_name = "${aws_elb.console.dns_name}"
    zone_id  = "${aws_elb.console.zone_id}"
  }
}
