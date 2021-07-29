############################################
# Datasources
############################################

data "ibm_cis_domain" "base_domain" {
  cis_id = var.cis_id
  domain = var.base_domain
}

############################################
# CIS DNS records (CNAME)
############################################

resource "ibm_cis_dns_record" "kubernetes_api" {
  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.base_domain.id
  type      = "CNAME"
  name      = "api.${var.cluster_domain}"
  content   = var.lb_kubernetes_api_public_hostname != "" ? var.lb_kubernetes_api_public_hostname : var.lb_kubernetes_api_private_hostname
  ttl       = 60
}

resource "ibm_cis_dns_record" "kubernetes_api_internal" {
  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.base_domain.id
  type      = "CNAME"
  name      = "api-int.${var.cluster_domain}"
  content   = var.lb_kubernetes_api_private_hostname
  ttl       = 60
}
