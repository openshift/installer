# ================COMMON=====================

variable "cluster_id" {
  type = string

  description = <<EOF
(internal) This is an identifier that can uniquely identify the cluster.

All the resources created include `cluster_id` for uniquness purposes.
EOF

}

variable "cluster_domain" {
  type = string

  description = <<EOF
The domain of the cluster.
All the records for the cluster are created under this domain.
Note: This field MUST be set manually prior to creating the cluster.
EOF

}

variable "bootstrap_ign_file" {
  type = string

  description = <<EOF
The file that contains the Ignition config used to configure the RHCOS based bootstrap machine.
EOF

}

variable "master_ign_file" {
  type = string

  description = <<EOF
The file that contains the Ignition config used to configure the RHCOS based control plane machines.
EOF

}

variable "worker_ign_file" {
  type = string

  description = <<EOF
The file that contains the Ignition config used to configure the RHCOS based worker machines.
EOF

}

variable "master_count" {
  type    = string
  default = "1"

  description = <<EOF
The number of control plane machines required.

Since etcd is colocated on control plane machines, suggested number is 3 or 5.
Default: 1
EOF

}

variable "worker_count" {
  type    = string
  default = "1"

  description = <<EOF
The number of worker machines required.

Default: 1
EOF

}

# ================MATCHBOX=====================

variable "matchbox_rpc_endpoint" {
  type = string

  description = <<EOF
RPC endpoint for matchbox.

For more info: https://godoc.org/github.com/coreos/matchbox/matchbox/client
EOF

}

variable "matchbox_http_endpoint" {
  type = string

  description = <<EOF
HTTPS endpoint for matchbox. This must include the scheme

For more info: https://github.com/coreos/matchbox/blob/master/Documentation/api.md
EOF

}

variable "matchbox_trusted_ca_cert" {
  type    = string
  default = "matchbox/tls/ca.crt"

  description = <<EOF
Certificate Authority certificate to trust the matchbox endpoint.
EOF

}

variable "matchbox_client_cert" {
  type    = string
  default = "matchbox/tls/client.crt"

  description = <<EOF
Client certificate used to authenticate with the matchbox RPC API.

For more info: https://github.com/coreos/matchbox/blob/master/Documentation/api.md
EOF

}

variable "matchbox_client_key" {
  type    = string
  default = "matchbox/tls/client.key"

  description = <<EOF
Client certificate's key used to authenticate with the matchbox RPC API.

For more info: https://github.com/coreos/matchbox/blob/master/Documentation/api.md
EOF

}

variable "pxe_kernel_args" {
  type    = string
  default = ""

  description = <<EOF
Arbitrary kernel arguments, space delimited ie:
coreos.inst.image_url=http://example.com/image.gz coreos.color=blue
EOF

}

variable "pxe_kernel_url" {
  type = string

  description = <<EOF
URL to the kernel image that should be used to PXE machines.

This can be a fully-qualified URL or URL relative to matchbox_http_endpoint to use Matchbox assets (https://github.com/coreos/matchbox/blob/master/Documentation/matchbox.md#assets).
EOF

}

variable "pxe_initrd_url" {
  type = string

  description = <<EOF
URL to the initrd image that should be used to PXE machines.

This can be a fully-qualified URL or URL relative to matchbox_http_endpoint to use Matchbox assets (https://github.com/coreos/matchbox/blob/master/Documentation/matchbox.md#assets).
EOF

}

# ================METAL=====================

variable "metal_project_id" {
  type = string

  description = <<EOF
The Project ID for Equinix Metal where servers will be deployed.
EOF

}

variable "metal_plan" {
  type    = string
  default = "c1.small.x86"

  description = <<EOF
The Equinix Metal device plan slug.
EOF

}

variable "metal_facility" {
  type    = string
  default = "any"

  description = <<EOF
The Equinix Metal facilities code to be used.
EOF

}

variable "metal_hardware_reservation_id" {
  type    = string
  default = ""

  description = <<EOF
The UUID of the hardware reservation where you want this device deployed on Equinix Metal.
EOF

}

# ================AWS=====================

variable "public_r53_zone" {
  type = string

  description = <<EOF
The name of the public route53 zone that should be used to create DNS records for the cluster.
EOF

}

variable "bootstrap_dns" {
  default = true

  description = <<EOF
(internal) This defines if the bootstrap machine should be part of the API pool.

Default: true
EOF

}
