variable "master_count" {
  type = "string"
}

variable "worker_count" {
  type = "string"
}

// The name of the cluster.
variable "cluster_name" {
  type = "string"
}

variable "external_gateway_id" {
  type = "string"
}

variable "service_cidr" {
  type = "string"
}

variable "cluster_cidr" {
  type = "string"
}

variable "floatingip_pool" {
  type    = "string"
  default = "public"
}

variable "subnet_cidr" {
  type    = "string"
  default = "192.168.1.0/24"
}

variable "dns_nameservers" {
  type    = "list"
  default = ["8.8.8.8", "8.8.4.4"]
}

output "master_ports" {
  value = ["${openstack_networking_port_v2.master.*.id}"]
}

output "master_floating_ips" {
  value = ["${openstack_networking_floatingip_v2.master.*.address}"]
}

output "worker_ports" {
  value = ["${openstack_networking_port_v2.worker.*.id}"]
}

output "worker_floating_ips" {
  value = ["${openstack_networking_floatingip_v2.worker.*.address}"]
}

output "id" {
  value = "${openstack_networking_network_v2.network.id}"
}
