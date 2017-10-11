output "version" {
  value = "${var.version == "latest" ? data.external.version.result["version"] : var.version}"
}
