resource "azurerm_virtual_machine" "etcd_node" {
  count                 = "${var.tectonic_etcd_count}"
  name                  = "${var.tectonic_cluster_name}_etcd_node_${count.index}"
  resource_group_name   = "${var.tectonic_azure_resource_group_name}"
  network_interface_ids = ["${azurerm_network_interface.etcd_nic.id}"]
  vm_size               = "${var.tectonic_azure_vm_size}"
  location              = "${var.tectonic_azure_location}"

  storage_image_reference {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "latest"
  }

  storage_os_disk {
    name          = "etcd-disk"
    vhd_uri       = "${azurerm_storage_account.etcd_storage.primary_blob_endpoint}${azurerm_storage_container.etcd_storage_container.name}/etcd-disk.vhd"
    caching       = "ReadWrite"
    create_option = "FromImage"
  }

  os_profile {
    computer_name  = "etcd"
    admin_username = "core"
    admin_password = "Microsoft123!"
    custom_data    = "${base64encode("${ignition_config.etcd.rendered}")}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_storage_account" "etcd_storage" {
  name                = "${var.tectonic_cluster_name}etcdstorage2"
  resource_group_name = "${var.tectonic_azure_resource_group_name}"
  location            = "${var.tectonic_azure_location}"
  account_type        = "Standard_LRS"
}

resource "azurerm_storage_container" "etcd_storage_container" {
  name                  = "${var.tectonic_cluster_name}-etcd-storage-container"
  resource_group_name   = "${var.tectonic_azure_resource_group_name}"
  storage_account_name  = "${azurerm_storage_account.etcd_storage.name}"
  container_access_type = "private"
  depends_on            = ["azurerm_storage_account.etcd_storage"]
}

resource "azurerm_network_interface" "etcd_nic" {
  name                      = "${var.tectonic_cluster_name}_etcd_nic"
  location                  = "${var.tectonic_azure_location}"
  network_security_group_id = "${azurerm_network_security_group.etcd_group.id}"
  resource_group_name       = "${var.tectonic_azure_resource_group_name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.etcd_subnet.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.etcd_publicip.id}"
  }
}

# resource "azurerm_dns_a_record" "tectonic-etcd" {
#   resource_group_name = "${var.tectonic_azure_resource_group_name}"
#   zone_name           = "${var.tectonic_azure_etc_dns_zone_name}"

#   name    = "${var.tectonic_cluster_name}-etcd"
#   ttl     = "60"
#   records = ["${azurerm_public_ip.etcd_publicip.ip_address}"]
# }

resource "azurerm_public_ip" "etcd_publicip" {
  name                         = "${var.tectonic_cluster_name}_etcd_publicip"
  location                     = "${var.tectonic_azure_location}"
  resource_group_name          = "${var.tectonic_azure_resource_group_name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_subnet" "etcd_subnet" {
  name                 = "${var.tectonic_cluster_name}_etcd_subnet"
  resource_group_name  = "${var.tectonic_azure_resource_group_name}"
  virtual_network_name = "${azurerm_virtual_network.etcd_vnet.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_network" "etcd_vnet" {
  name                = "${var.tectonic_cluster_name}_etcd_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${var.tectonic_azure_resource_group_name}"
}

resource "azurerm_network_security_group" "etcd_group" {
  name                = "${var.tectonic_cluster_name}_etcd_group"
  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${var.tectonic_azure_resource_group_name}"

  security_rule {
    name                       = "rule1"
    source_port_range          = 22
    destination_port_range     = 22
    protocol                   = "Tcp"
    destination_address_prefix = "0.0.0.0/0"
    source_address_prefix      = "0.0.0.0/0"
    access                     = "Allow"
    priority                   = "100"
    direction                  = "Inbound"
  }

  security_rule {
    name                       = "rule2"
    source_port_range          = 2379
    destination_port_range     = 2380
    protocol                   = "Tcp"
    destination_address_prefix = "0.0.0.0/0"
    source_address_prefix      = "0.0.0.0/0"
    access                     = "Allow"
    priority                   = "101"
    direction                  = "Inbound"
  }

  security_rule {
    name                       = "rule3"
    source_port_range          = "*"
    destination_port_range     = "*"
    protocol                   = "*"
    destination_address_prefix = "0.0.0.0/0"
    source_address_prefix      = "0.0.0.0/0"
    access                     = "Allow"
    priority                   = "103"
    direction                  = "Inbound"
  }

  security_rule {
    name                       = "rule4"
    source_port_range          = "*"
    destination_port_range     = "*"
    protocol                   = "*"
    destination_address_prefix = "Internet"
    source_address_prefix      = "0.0.0.0/0"
    access                     = "Allow"
    priority                   = "104"
    direction                  = "Outbound"
  }
}
