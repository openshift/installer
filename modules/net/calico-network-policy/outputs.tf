output "id" {
  value = "${ var.enabled ? "${sha1("${join(" ", local_file.calico-network-policy.*.id)}")}" : "# calico policy disabled" }"
}

output "name" {
  value = "calico"
}
