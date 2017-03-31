// The flavor ID as given in `openstack flavor list`.
// Specifies the size (CPU/Memory/Drive) of the VM.
variable "tectonic_openstack_flavor_id" {
  type    = "string"
  default = "5cf64088-893b-46b5-9bb1-ee020277635d"
}

// The image ID as given in `openstack image list`.
// Specifies the OS image of the VM.
variable "tectonic_openstack_image_id" {
  type    = "string"
  default = "acdcd535-5408-40f3-8e88-ad8ebb6507e6"
}

// The ID of the network to be used as the external internet gateway
// as given in `openstack network list`.
variable "tectonic_openstack_external_gateway_id" {
  type    = "string"
  default = "6d6357ac-0f70-4afa-8bd7-c274cc4ea235"
}

// The name name of the floating IP pool
// as given in `openstack floating ip list`.
// This pool will be used to assign floating IPs to worker and master nodes.
variable "tectonic_openstack_floatingip_pool" {
  type    = "string"
  default = "public"
}

// The subnet CIDR for the master/worker/etcd compute nodes.
// This CIDR will also be assigned to the created the OpenStack subnet resource.
variable "tectonic_openstack_subnet_cidr" {
  type    = "string"
  default = "192.168.1.0/24"
}

// The DNS servers assigned to the generated OpenStack subnet resource.
variable "tectonic_openstack_dns_nameservers" {
  type    = "list"
  default = ["8.8.8.8", "8.8.4.4"]
}
