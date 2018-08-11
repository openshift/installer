output "etcd_count" {
  value = "${var.etcd_count > 0 ? var.etcd_count : 1}"

  description = <<EOF
The number of etcd nodes to be created.
This will always be greater than zero.
EOF
}
