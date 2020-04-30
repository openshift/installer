variable "base_image_id" {
  type        = string
  description = "The identifier of the Glance image for the bootstrap node."
}

variable "extra_tags" {
  type    = map(string)
  default = {}

  description = <<EOF
(optional) Extra tags to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

}

variable "cluster_id" {
  type = string
  description = "The identifier for the cluster."
}

variable "ignition" {
  type = string
  description = "The content of the bootstrap ignition file."
}

variable "bootstrap_shim_ignition" {
  type = string
  description = "The content of the ignition file with user ca bundle."
}

variable "flavor_name" {
  type = string
  description = "The Nova flavor for the bootstrap node."
}

variable "api_int_ip" {
  type = string
}

variable "node_dns_ip" {
  type = string
}

variable "external_network" {
  type = string
}

variable "private_network_id" {
  type = string
}

variable "master_sg_id" {
  type = string
}

variable "nodes_subnet_id" {
  type = string
}

variable "cluster_domain" {
  type = string
}

variable "master_port_ids" {
  type = list(string)
}

variable "root_volume_size" {
  type = number
  description = "The size of the volume in gigabytes for the root block device."
}

variable "root_volume_type" {
  type = string
  description = "The type of volume for the root block device."
}
