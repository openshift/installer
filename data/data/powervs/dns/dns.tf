data "ibm_dns_domain" "domain" {
  name = var.base_domain
}
resource "ibm_dns_record" "api" {
  data               = "${var.load_balancer_hostname}."
  domain_id          = data.ibm_dns_domain.domain.id
  host               = "api.${var.cluster_domain}"
  responsible_person = "root.${var.cluster_domain}."
  ttl                = 900
  type               = "cname"
  tags               = [var.cluster_id, "${var.cluster_id}-api"]
}
resource "ibm_dns_record" "api-int" {
  data               = "${var.load_balancer_int_hostname}."
  domain_id          = data.ibm_dns_domain.domain.id
  host               = "api-int.${var.cluster_domain}"
  responsible_person = "root.${var.cluster_domain}."
  ttl                = 900
  type               = "cname"
  tags               = [var.cluster_id, "${var.cluster_id}-api-int"]
}
resource "ibm_dns_record" "apps" {
  data               = "${var.load_balancer_hostname}."
  domain_id          = data.ibm_dns_domain.domain.id
  host               = "*.apps.${var.cluster_domain}"
  responsible_person = "root.${var.cluster_domain}."
  ttl                = 900
  type               = "cname"
  tags               = [var.cluster_id, "${var.cluster_id}-apps"]
}
