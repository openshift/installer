resource "azurerm_virtual_machine" "etcd_node" {
  count                 = "${var.tectonic_etcd_count}"
  name                  = "${var.tectonic_cluster_name}_etcd_node_${count.index}"
  resource_group_name   = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  network_interface_ids = ["${azurerm_network_interface.etcd_nic.id}"]
  vm_size               = "${var.tectonic_azure_vm_size}"
  location              = "East US"

  metadata {
    role = "etcd"
  }

  storage_image_reference {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "latest"
  }

  os_profile {
    computer_name  = "etcd"
    admin_username = "${var.admin_username}"
    admin_password = "microsoft123"
    custom_data    = "${base64encode(file("${path.module}/userdata.yml"))}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_network_interface" "etcd_nic" {
  name                      = "${var.tectonic_cluster_name}_etcd_nic"
  location                  = "East US"
  network_security_group_id = "${azurerm_network_security_group.etcd_group.id}"
  resource_group_name       = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.etcd_subnet.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.etcd_publicip.id}"
  }
}

resource "azurerm_public_ip" "etcd_publicip" {
  name                         = "${var.tectonic_cluster_name}_etcd_publicip"
  location                     = "East US"
  resource_group_name          = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_subnet" "etcd_subnet" {
  name                 = "${var.tectonic_cluster_name}_etcd_subnet"
  resource_group_name  = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  virtual_network_name = "${azurerm_virtual_network.etcd_vnet.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_network" "etcd_vnet" {
  name                = "${var.tectonic_cluster_name}_etcd_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "East US"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
}

resource "azurerm_network_security_group" "etcd_group" {
  name = "${var.tectonic_cluster_name}_etcd_group"

  security_rule {
    source_port_range      = 22
    destination_port_range = 22
    protocol               = "tcp"
    source_address_prefix  = "0.0.0.0/0"
  }

  security_rule {
    source_port_range      = 2379
    destination_port_range = 2380
    protocol               = "tcp"
    source_address_prefix  = "0.0.0.0/0"
  }

  security_rule {
    source_port_range      = -1
    destination_port_range = -1
    protocol               = "icmp"
    source_address_prefix  = "0.0.0.0/0"
  }
}
