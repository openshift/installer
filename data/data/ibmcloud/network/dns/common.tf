locals {
  dns_zone_id = var.is_external ? "" : data.ibm_dns_zones.zones[0].dns_zones[index(data.ibm_dns_zones.zones[0].dns_zones[*].name, var.base_domain)].zone_id
}

############################################
# DNS Zone
############################################

data "ibm_dns_zones" "zones" {
  count = var.is_external ? 0 : 1

  instance_id = var.dns_id
}
