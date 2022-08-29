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

locals {
  proxy_count = var.publish_strategy == "Internal" ? 1 : 0
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

data "ibm_is_image" "dns_vm_image" {
  count = local.proxy_count
  name  = var.dns_vm_image_name
}

locals {
  user_data_string = <<EOF
#cloud-config
packages:
  - bind
  - bind-utils
runcmd:
  - systemctl enable named.service
  - systemctl start named.service
EOF
}

#
# The following is because ci/prow/tf-fmt is recommending that
# style of formatting which seems like a bug to me.
#
resource "ibm_is_instance" "dns_vm_vsi" {
  count = local.proxy_count
  name = "${var.cluster_id}-dns-vsi"
  vpc = var.vpc_id
  zone = var.vpc_zone
  keys = [ibm_is_ssh_key.dns_ssh_key[0].id]
  image = data.ibm_is_image.dns_vm_image[0].id
  profile = "cx2-2x4"

  primary_network_interface {
    subnet = var.vpc_subnet_id
    security_groups = [ibm_is_security_group.dns_vm_sg[0].id]
  }

  user_data = base64encode(local.user_data_string)
}
