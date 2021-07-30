data "ibm_cis_domain" "base_domain" {
  cis_id = var.cis_id
  domain = var.base_domain
}

resource "ibm_cis_dns_record" "kubernetes_api" {
  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.base_domain.id
  type      = "CNAME"
  name      = "api.${var.cluster_domain}"
  content   = var.load_balancer_hostname
  ttl       = 60
}

resource "ibm_cis_dns_record" "kubernetes_api_internal" {
  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.base_domain.id
  type      = "CNAME"
  name      = "api-int.${var.cluster_domain}"
  content   = var.load_balancer_int_hostname
  ttl       = 60
}

resource "ibm_cis_dns_record" "apps" {
  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.base_domain.id
  type      = "CNAME"
  name      = "*.apps.${var.cluster_domain}"
  content   = var.load_balancer_hostname
  ttl       = 60
}
