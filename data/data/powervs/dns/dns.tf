locals {
  // extracting "cluster_name" from <cluster_domain> := <cluster_name>.<base_domain>
  cluster_name = replace(var.cluster_domain, ".${var.base_domain}", "")
}

data "ibm_dns_domain" "domain" {
  name = var.base_domain
}
resource "ibm_dns_record" "api" {
  data               = "${var.load_balancer_hostname}."
  domain_id          = data.ibm_dns_domain.domain.id
  host               = "api.${local.cluster_name}"
  responsible_person = "root.${local.cluster_name}."
  ttl                = 900
  type               = "cname"
  tags               = [var.cluster_id, "${local.cluster_name}-api"]
}
resource "ibm_dns_record" "api-int" {
  data               = "${var.load_balancer_int_hostname}."
  domain_id          = data.ibm_dns_domain.domain.id
  host               = "api-int.${local.cluster_name}"
  responsible_person = "root.${local.cluster_name}."
  ttl                = 900
  type               = "cname"
  tags               = [var.cluster_id, "${local.cluster_name}-api-int"]
}
resource "ibm_dns_record" "apps" {
  data               = "${var.load_balancer_hostname}."
  domain_id          = data.ibm_dns_domain.domain.id
  host               = "*.apps.${local.cluster_name}"
  responsible_person = "root.${local.cluster_name}."
  ttl                = 900
  type               = "cname"
  tags               = [var.cluster_id, "${local.cluster_name}-apps"]
}
