output "etcd_count" {
  value = "${var.etcd_count > 0 ? var.etcd_count : max(local.zone_count_odd, 3)}"

  description = <<EOF
The number of etcd nodes to be created.
This will always be greater than zero.
EOF
}
