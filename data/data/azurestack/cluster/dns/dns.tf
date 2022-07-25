locals {
  // extracting <clustername> from <clusterdomain>
  cluster_name = replace(var.cluster_domain, ".${var.base_domain}", "")
}

resource "azurestack_dns_a_record" "api_external_v4" {
  name                = "api.${local.cluster_name}"
  zone_name           = var.base_domain
  resource_group_name = var.base_domain_resource_group_name
  ttl                 = 300
  records             = var.private ? [var.ilb_ipaddress_v4] : [var.elb_pip_v4]
  tags                = var.tags
}

resource "azurestack_dns_a_record" "api_internal_v4" {
  name                = "api-int.${local.cluster_name}"
  zone_name           = var.base_domain
  resource_group_name = var.base_domain_resource_group_name
  ttl                 = 300
  records             = [var.ilb_ipaddress_v4]
  tags                = var.tags
}
