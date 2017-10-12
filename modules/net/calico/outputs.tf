output "id" {
  value = "${var.enabled ? "${sha1("${join(" ", local_file.calico.*.id)}")}" : "# calico disabled"}"
}

output "name" {
  value = "calico"
}
