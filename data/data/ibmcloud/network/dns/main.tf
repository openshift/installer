############################################
# DNS permitted networks
############################################

resource "ibm_dns_permitted_network" "vpc" {
  # Only create the Permitted Network if Internal (Private using DNS) and the VPC is not already a Permitted Network
  count = ! var.is_external && ! var.vpc_permitted ? 1 : 0

  instance_id = var.dns_id
  zone_id     = local.dns_zone_id
  vpc_crn     = var.vpc_crn
  type        = "vpc"
}

############################################
# DNS records (CNAME)
############################################

resource "ibm_dns_resource_record" "kubernetes_api_internal_public" {
  count = var.is_external ? 0 : 1

  instance_id = var.dns_id
  zone_id     = local.dns_zone_id
  type        = "CNAME"
  name        = "api.${var.cluster_domain}"
  rdata       = var.lb_kubernetes_api_private_hostname
  ttl         = "60"
}

resource "ibm_dns_resource_record" "kubernetes_api_private" {
  count = var.is_external ? 0 : 1

  instance_id = var.dns_id
  zone_id     = local.dns_zone_id
  type        = "CNAME"
  name        = "api-int.${var.cluster_domain}"
  rdata       = var.lb_kubernetes_api_private_hostname
  ttl         = "60"
}
