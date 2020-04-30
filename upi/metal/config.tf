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

variable "pxe_os_image_url" {
  type = string

  description = <<EOF
URL to the OS image for RHCOS that should be installed on machines.

For more info: https://github.com/coreos/coreos-installer#kernel-command-line-options-for-coreos-installer-running-in-the-initramfs
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

# ================PACKET=====================

variable "packet_project_id" {
  type = string

  description = <<EOF
The Project ID for Packet.net where servers will be deployed.
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
