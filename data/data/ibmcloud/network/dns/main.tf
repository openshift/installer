############################################
# DNS permitted networks
############################################

resource "ibm_dns_permitted_network" "vpc" {
  count = var.is_external ? 0 : 1

  instance_id = var.dns_id
  zone_id     = local.dns_zone_id
  vpc_crn     = var.vpc_crn
  type        = "vpc"
}

############################################
# DNS records (CNAME)
############################################

resource "ibm_dns_resource_record" "kubernetes_api_private" {
  count = var.is_external ? 0 : 1

  instance_id = var.dns_id
  zone_id     = local.dns_zone_id
  type        = "CNAME"
  name        = "api-int.${var.cluster_domain}"
  rdata       = var.lb_kubernetes_api_private_hostname
  ttl         = "60"
}
