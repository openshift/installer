variable "tectonic_openstack_nova_config_version" {
  description = <<EOF
(internal) This declares the version of the OpenStack Nova configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF

  default = "1.0"
}

variable "tectonic_openstack_flavor_id" {
  type = "string"

  description = <<EOF
The flavor ID as given in `openstack flavor list`.
Specifies the size (CPU/Memory/Drive) of the VM.
EOF
}

variable "tectonic_openstack_image_id" {
  type = "string"

  description = <<EOF
The image ID as given in `openstack image list`.
Specifies the OS image of the VM.
EOF
}

variable "tectonic_openstack_network_name" {
  type = "string"

  description = <<EOF
The name of the network to connect to OpenStack compute instances.
EOF

  default = "public"
}
