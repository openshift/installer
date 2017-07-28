output "id" {
  value = "${ var.enabled ? "${sha1("${join(" ", local_file.calico_network_policy.*.id)}")}" : "# calico policy disabled" }"
}

output "name" {
  value = "calico"
}
