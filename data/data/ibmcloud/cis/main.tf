locals {
  prefix = var.cluster_id
}

############################################
# Datasources
############################################

data "ibm_cis_domain" "domain" {
  cis_id = var.cis_id
  domain = var.domain
}

############################################
# CIS DNS records (CNAME)
############################################

resource "ibm_cis_dns_record" "kubernetes_api_public" {
  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.domain.id
  type      = "CNAME"
  name      = "api.${var.cluster_id}"
  content   = var.lb_kubernetes_api_public_hostname
  ttl       = 60
}

resource "ibm_cis_dns_record" "kubernetes_api_private" {
  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.domain.id
  type      = "CNAME"
  name      = "api-int.${var.cluster_id}"
  content   = var.lb_kubernetes_api_private_hostname
  ttl       = 60
}

############################################
# CIS DNS records (A)
############################################

resource "ibm_cis_dns_record" "bootstrap_node" {
  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.domain.id
  type      = "A"
  name      = var.bootstrap_name
  content   = var.bootstrap_ipv4_address
  ttl       = 60
}

resource "ibm_cis_dns_record" "master_node" {
  count     = var.master_count

  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.domain.id
  type      = "A"
  name      = var.master_name_list[count.index]
  content   = var.master_ipv4_address_list[count.index]
  ttl       = 60
}

############################################
# CIS DNS records (PTR)
############################################

resource "ibm_cis_dns_record" "bootstrap_node_ptr" {
  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.domain.id
  type      = "PTR"
  name      = var.bootstrap_ipv4_address
  content   = "${var.bootstrap_name}.${var.cluster_id}.${var.domain}"
  ttl       = 60
}

resource "ibm_cis_dns_record" "master_node_ptr" {
  count     = var.master_count

  cis_id    = var.cis_id
  domain_id = data.ibm_cis_domain.domain.id
  type      = "PTR"
  name      = var.master_ipv4_address_list[count.index]
  content   = "${var.master_name_list[count.index]}.${var.cluster_id}.${var.domain}"
  ttl       = 60
}
