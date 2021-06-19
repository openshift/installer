locals {
  tags = merge(
    {
      "kubernetes.io_cluster.${var.cluster_id}" = "owned"
    },
    var.azure_extra_tags,
  )
  description = "Created By OpenShift Installer"
}

provider "azurestack" {
  arm_endpoint    = var.azure_arm_endpoint
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}

module "bootstrap" {
  source                 = "./bootstrap"
  resource_group_name    = data.azurestack_resource_group.main.name
  region                 = var.azure_region
  vm_size                = var.azure_bootstrap_vm_type
  vm_image_uri           = azurestack_image.cluster.id
  cluster_id             = var.cluster_id
  ignition               = var.ignition_bootstrap
  subnet_id              = module.vnet.master_subnet_id
  elb_backend_pool_v4_id = module.vnet.public_lb_backend_pool_v4_id
  ilb_backend_pool_v4_id = module.vnet.internal_lb_backend_pool_v4_id
  tags                   = local.tags
  storage_account        = azurestack_storage_account.cluster
  nsg_name               = module.vnet.cluster_nsg_name
  private                = var.azure_private
  availability_set_id    = azurestack_availability_set.master_availability_set.id
}

module "vnet" {
  source              = "./vnet"
  resource_group_name = data.azurestack_resource_group.main.name
  vnet_v4_cidrs       = var.machine_v4_cidrs
  cluster_id          = var.cluster_id
  region              = var.azure_region
  dns_label           = var.cluster_id

  preexisting_network         = var.azure_preexisting_network
  network_resource_group_name = var.azure_network_resource_group_name
  virtual_network_name        = var.azure_virtual_network
  master_subnet               = var.azure_control_plane_subnet
  worker_subnet               = var.azure_compute_subnet
  private                     = var.azure_private
}

module "master" {
  source                 = "./master"
  resource_group_name    = data.azurestack_resource_group.main.name
  cluster_id             = var.cluster_id
  region                 = var.azure_region
  vm_size                = var.azure_master_vm_type
  vm_image_uri           = azurestack_image.cluster.id
  ignition               = var.ignition_master
  elb_backend_pool_v4_id = module.vnet.public_lb_backend_pool_v4_id
  ilb_backend_pool_v4_id = module.vnet.internal_lb_backend_pool_v4_id
  subnet_id              = module.vnet.master_subnet_id
  instance_count         = var.master_count
  storage_account        = azurestack_storage_account.cluster
  os_volume_size         = var.azure_master_root_volume_size
  private                = var.azure_private
  availability_set_id    = azurestack_availability_set.master_availability_set.id
}

module "dns" {
  source                          = "./dns"
  cluster_domain                  = var.cluster_domain
  cluster_id                      = var.cluster_id
  base_domain                     = var.base_domain
  virtual_network_id              = module.vnet.virtual_network_id
  external_lb_fqdn_v4             = module.vnet.public_lb_pip_v4_fqdn
  external_lb_pip_v4              = module.vnet.public_lb_pip_v4
  internal_lb_ipaddress_v4        = module.vnet.internal_lb_ip_v4_address
  resource_group_name             = data.azurestack_resource_group.main.name
  base_domain_resource_group_name = var.azure_base_domain_resource_group_name
  private                         = var.azure_private
}

resource "random_string" "storage_suffix" {
  length  = 5
  upper   = false
  special = false
}

resource "azurestack_resource_group" "main" {
  count = var.azure_resource_group_name == "" ? 1 : 0

  name     = "${var.cluster_id}-rg"
  location = var.azure_region
  tags     = local.tags
}

data "azurestack_resource_group" "main" {
  name = var.azure_resource_group_name == "" ? "${var.cluster_id}-rg" : var.azure_resource_group_name

  depends_on = [azurestack_resource_group.main]
}

data "azurestack_resource_group" "network" {
  count = var.azure_preexisting_network ? 1 : 0

  name = var.azure_network_resource_group_name
}

resource "azurestack_storage_account" "cluster" {
  name                     = "cluster${random_string.storage_suffix.result}"
  resource_group_name      = data.azurestack_resource_group.main.name
  location                 = var.azure_region
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# copy over the vhd to cluster resource group and create an image using that
resource "azurestack_storage_container" "vhd" {
  name                 = "vhd"
  resource_group_name  = data.azurestack_resource_group.main.name
  storage_account_name = azurestack_storage_account.cluster.name
}

resource "azurestack_storage_blob" "rhcos_image" {
  name                   = "rhcos${random_string.storage_suffix.result}.vhd"
  resource_group_name    = data.azurestack_resource_group.main.name
  storage_account_name   = azurestack_storage_account.cluster.name
  storage_container_name = azurestack_storage_container.vhd.name
  type                   = "page"
  source_uri             = var.azure_image_url
}

resource "azurestack_image" "cluster" {
  name                = var.cluster_id
  resource_group_name = data.azurestack_resource_group.main.name
  location            = var.azure_region

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = azurestack_storage_blob.rhcos_image.url
  }
}

resource "azurestack_availability_set" "master_availability_set" {
  name                = "${var.cluster_id}-master"
  resource_group_name = data.azurestack_resource_group.main.name
  location            = var.azure_region
  managed             = true
}
