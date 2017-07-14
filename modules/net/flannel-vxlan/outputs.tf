output "id" {
  value = "${ var.enabled ? "${sha1("${join(" ", local_file.flannel.*.id)}")}" : "# flannel disabled"}"
}

output "name" {
  value = "flannel-vxlan"
}
