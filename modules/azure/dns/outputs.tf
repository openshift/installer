output "ingress_external_fqdn" {
  value = "${var.cluster_name}.${azurerm_dns_zone.tectonic_azure_dns_zone.name}"
}

output "ingress_internal_fqdn" {
  value = "${var.cluster_name}.${azurerm_dns_zone.tectonic_azure_dns_zone.name}"
}

output "api_external_fqdn" {
  value = "${var.cluster_name}-k8s.${azurerm_dns_zone.tectonic_azure_dns_zone.name}"
}

output "api_internal_fqdn" {
  value = "${var.cluster_name}-k8s.${azurerm_dns_zone.tectonic_azure_dns_zone.name}"
}
