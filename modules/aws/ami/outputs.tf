output "id" {
  value       = "${data.aws_ami.coreos_ami.image_id}"
  description = "The selected CoreOS Container Linux AMI ID."
}
