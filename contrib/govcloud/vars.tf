# placeholders for access_key / secret_key
# should be fed through env var or variable file
# https://www.terraform.io/docs/configuration/variables.html

variable vpc_name {
  description = "The name of the VPC to identify created resources."
  default     = "bastion"
}

variable base_domain {
  default     = "tectonic-ci.de"
  description = "The base domain for this cluster's FQDN"
}

variable vpc_aws_region {
  description = "The target AWS region for the cluster"
  default     = "us-gov-west-1"
}

variable vpc_cidr {
  default     = "10.0.0.0/16"
  description = "The CIDR range used for your entire VPC"
}

variable subnet_count {
  default     = 4
  description = "Number of private subnets to pre-create"
}

variable local_network_cidr {
  default     = "10.7.0.0/16"
  description = "IP range in the network your laptop is on (dosn't actually matter unless your instances need to connect to the local network your laptop is on)"
}

variable "ssh_key" {
  default = ""
}

variable "nginx_username" {
  description = "Used for retrieving the OpenVPN client config file."
  default     = "user"
}

variable "nginx_password" {
  description = "Used for retrieving the OpenVPN client config file."
  default     = "password"
}

output "ovpn_url" {
  value = "http://${aws_eip.ovpn_eip.public_ip}"
}

output "base_domain" {
  value = "${var.base_domain}"
}

output "vpc_id" {
  value = "${aws_vpc.vpc.id}"
}

output "vpc_dns" {
  value = "${aws_instance.bastion.private_ip}"
}

output "subnets" {
  value = "${aws_subnet.priv_subnet.*.id}"
}
