############################################
# Datasources
############################################

data "ibm_cis_domain" "base_domain" {
  count = var.is_external ? 1 : 0

  cis_id = var.cis_id
  domain = var.base_domain
}

############################################
# CIS DNS records (CNAME)
############################################

resource "ibm_cis_dns_record" "kubernetes_api" {
  count = var.is_external ? 1 : 0

  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.base_domain[0].id
  type      = "CNAME"
  name      = "api.${var.cluster_domain}"
  content   = var.lb_kubernetes_api_public_hostname != "" ? var.lb_kubernetes_api_public_hostname : var.lb_kubernetes_api_private_hostname
  ttl       = 60
}

resource "ibm_cis_dns_record" "kubernetes_api_internal" {
  count = var.is_external ? 1 : 0

  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.base_domain[0].id
  type      = "CNAME"
  name      = "api-int.${var.cluster_domain}"
  content   = var.lb_kubernetes_api_private_hostname
  ttl       = 60
}
