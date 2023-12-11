locals {
  description = "Created By OpenShift Installer"
  # At this time min_tls_version is only supported in the Public Cloud and US Government Cloud.
  environments_with_min_tls_version = ["public", "usgovernment"]
}

provider "azurerm" {
  features {}
  subscription_id             = var.azure_subscription_id
  client_id                   = var.azure_client_id
  client_secret               = var.azure_client_secret
  client_certificate_password = var.azure_certificate_password
  client_certificate_path     = var.azure_certificate_path
  tenant_id                   = var.azure_tenant_id
  use_msi                     = var.azure_use_msi
  environment                 = var.azure_environment
}

module "master" {
  source                     = "./master"
  resource_group_name        = var.resource_group_name
  cluster_id                 = var.cluster_id
  region                     = var.azure_region
  availability_zones         = var.azure_master_availability_zones
  vm_size                    = var.azure_master_vm_type
  disk_encryption_set_id     = var.azure_master_disk_encryption_set_id
  encryption_at_host_enabled = var.azure_master_encryption_at_host_enabled
  vm_image                   = var.vm_image
  identity                   = var.identity
  ignition                   = var.ignition_master
  elb_backend_pool_v4_id     = var.elb_backend_pool_v4_id
  elb_backend_pool_v6_id     = var.elb_backend_pool_v6_id
  ilb_backend_pool_v4_id     = var.ilb_backend_pool_v4_id
  ilb_backend_pool_v6_id     = var.ilb_backend_pool_v6_id
  subnet_id                  = var.master_subnet_id
  instance_count             = var.master_count
  os_volume_type             = var.azure_master_root_volume_type
  os_volume_size             = var.azure_master_root_volume_size
  private                    = var.azure_private
  outbound_type              = var.azure_outbound_routing_type
  ultra_ssd_enabled          = var.azure_control_plane_ultra_ssd_enabled
  vm_networking_type         = var.azure_control_plane_vm_networking_type
  azure_extra_tags           = var.azure_extra_tags
  use_marketplace_image      = var.azure_use_marketplace_image
  vm_image_has_plan          = var.azure_marketplace_image_has_plan
  vm_image_publisher         = var.azure_marketplace_image_publisher
  vm_image_offer             = var.azure_marketplace_image_offer
  vm_image_sku               = var.azure_marketplace_image_sku
  vm_image_version           = var.azure_marketplace_image_version

  security_encryption_type            = var.azure_master_security_encryption_type
  secure_vm_disk_encryption_set_id    = var.azure_master_secure_vm_disk_encryption_set_id
  secure_boot                         = var.azure_master_secure_boot
  virtualized_trusted_platform_module = var.azure_master_virtualized_trusted_platform_module

  use_ipv4 = var.use_ipv4
  use_ipv6 = var.use_ipv6
}

module "dns" {
  source                          = "./dns"
  cluster_domain                  = var.cluster_domain
  cluster_id                      = var.cluster_id
  base_domain                     = var.base_domain
  virtual_network_id              = var.virtual_network_id
  external_lb_fqdn_v4             = var.public_lb_pip_v4_fqdn
  external_lb_fqdn_v6             = var.public_lb_pip_v6_fqdn
  internal_lb_ipaddress_v4        = var.internal_lb_ip_v4_address
  internal_lb_ipaddress_v6        = var.internal_lb_ip_v6_address
  resource_group_name             = var.resource_group_name
  base_domain_resource_group_name = var.azure_base_domain_resource_group_name
  private                         = var.azure_private
  azure_extra_tags                = var.azure_extra_tags

  use_ipv4 = var.use_ipv4
  use_ipv6 = var.use_ipv6
}
