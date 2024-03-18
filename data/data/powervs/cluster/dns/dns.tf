############
# Public DNS
############

data "ibm_cis_domain" "base_domain" {
  count  = var.publish_strategy == "Internal" ? 0 : 1
  cis_id = var.service_id
  domain = var.base_domain
}

resource "ibm_cis_dns_record" "kubernetes_api" {
  count     = var.publish_strategy == "Internal" ? 0 : 1
  cis_id    = var.service_id
  domain_id = data.ibm_cis_domain.base_domain[count.index].id
  type      = "CNAME"
  name      = "api.${var.cluster_domain}"
  content   = var.load_balancer_hostname
  ttl       = 60
}

resource "ibm_cis_dns_record" "kubernetes_api_internal" {
  count     = var.publish_strategy == "Internal" ? 0 : 1
  cis_id    = var.service_id
  domain_id = data.ibm_cis_domain.base_domain[count.index].id
  type      = "CNAME"
  name      = "api-int.${var.cluster_domain}"
  content   = var.load_balancer_int_hostname
  ttl       = 60
}

#############
# Private DNS
#############

data "ibm_dns_zones" "dns_zones" {
  count       = var.publish_strategy == "Internal" ? 1 : 0
  instance_id = var.service_id
}

resource "ibm_dns_permitted_network" "permit_vpc_network_for_dns" {
  count       = var.publish_strategy == "Internal" && ! var.vpc_permitted ? 1 : 0
  instance_id = var.service_id
  zone_id     = local.dns_zone.zone_id
  vpc_crn     = var.vpc_crn
  type        = "vpc"
}

resource "ibm_dns_resource_record" "kubernetes_api" {
  count       = var.publish_strategy == "Internal" ? 1 : 0
  instance_id = var.service_id
  zone_id     = local.dns_zone.zone_id
  type        = "CNAME"
  name        = "api.${var.cluster_domain}"
  rdata       = var.load_balancer_int_hostname
  ttl         = 60
}

resource "ibm_dns_resource_record" "kubernetes_api_internal" {
  count       = var.publish_strategy == "Internal" ? 1 : 0
  instance_id = var.service_id
  zone_id     = local.dns_zone.zone_id
  type        = "CNAME"
  name        = "api-int.${var.cluster_domain}"
  rdata       = var.load_balancer_int_hostname
  ttl         = 60
}

resource "ibm_dns_resource_record" "proxy_vsi_record" {
  count       = var.publish_strategy == "Internal" ? 1 : 0
  instance_id = var.service_id
  zone_id     = local.dns_zone.zone_id
  type        = "A"
  name        = "proxy.${var.cluster_domain}"
  rdata       = ibm_is_instance.dns_vm_vsi[0].primary_network_interface[0].primary_ip[0].address
  ttl         = 60
}

resource "ibm_is_ssh_key" "dns_ssh_key" {
  count      = local.proxy_count
  name       = "${var.cluster_id}-dns-ssh-key"
  public_key = var.ssh_key
}

resource "ibm_is_security_group" "dns_vm_sg" {
  count = local.proxy_count
  name  = "${var.cluster_id}-dns-sg"
  vpc   = var.vpc_id
}

# allow all outgoing network traffic
resource "ibm_is_security_group_rule" "dns_vm_sg_outgoing_all" {
  count     = local.proxy_count
  group     = ibm_is_security_group.dns_vm_sg[0].id
  direction = "outbound"
  remote    = "0.0.0.0/0"
}

# allow all incoming network traffic on port 22
resource "ibm_is_security_group_rule" "dns_vm_sg_ssh_all" {
  count     = local.proxy_count
  group     = ibm_is_security_group.dns_vm_sg[0].id
  direction = "inbound"
  remote    = "0.0.0.0/0"

  tcp {
    port_min = 22
    port_max = 22
  }
}

# allow all incoming network traffic on port 53
resource "ibm_is_security_group_rule" "dns_vm_sg_dns_all" {
  count     = local.proxy_count
  group     = ibm_is_security_group.dns_vm_sg[0].id
  direction = "inbound"
  remote    = "0.0.0.0/0"

  udp {
    port_min = 53
    port_max = 53
  }
}

# allow all incoming network traffic on port 80
resource "ibm_is_security_group_rule" "dns_vm_sg_http_all" {
  count     = local.proxy_count
  group     = ibm_is_security_group.dns_vm_sg[0].id
  direction = "inbound"
  remote    = "0.0.0.0/0"

  tcp {
    port_min = 80
    port_max = 80
  }
}

# allow all incoming network traffic on port 3128 for squid proxy
resource "ibm_is_security_group_rule" "dns_vm_sg_squid_all" {
  count     = local.proxy_count
  group     = ibm_is_security_group.dns_vm_sg[0].id
  direction = "inbound"
  remote    = "0.0.0.0/0"

  tcp {
    port_min = 3128
    port_max = 3128
  }
}

data "ibm_is_images" "images" {
  count  = local.proxy_count
  status = "available"
}

data "ibm_is_image" "dns_vm_image" {
  count = local.proxy_count
  name  = [for image in data.ibm_is_images.images[0].images : image if startswith(image.os, var.dns_vm_image_os)][0].name
}

locals {
  dns_zone    = var.publish_strategy == "Internal" ? data.ibm_dns_zones.dns_zones[0].dns_zones[index(data.ibm_dns_zones.dns_zones[0].dns_zones.*.name, var.base_domain)] : null
  proxy_count = var.publish_strategy == "Internal" ? 1 : 0
}

resource "ibm_is_instance" "dns_vm_vsi" {
  count   = local.proxy_count
  name    = "${var.cluster_id}-dns-vsi"
  vpc     = var.vpc_id
  zone    = var.vpc_zone
  keys    = [ibm_is_ssh_key.dns_ssh_key[0].id]
  image   = data.ibm_is_image.dns_vm_image[0].id
  profile = "cx2-2x4"

  primary_network_interface {
    subnet          = var.vpc_subnet_id
    security_groups = [ibm_is_security_group.dns_vm_sg[0].id]
  }

  user_data = templatefile("${path.module}/templates/cloud-init.yaml.tpl", { is_proxy : ! var.enable_snat, vpc_region : var.vpc_region })
}
