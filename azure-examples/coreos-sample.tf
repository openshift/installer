variable "resourcesname" {
  default = "jzcoreosterraform"
}

variable "azurelocation" {
    default = "EastUS"
}

variable "admin_username" {
    default = "jimzim"
}

# Configure the Microsoft Azure Provider
provider "azurerm" {
  subscription_id = ""
  client_id       = ""
  client_secret   = ""
  tenant_id       = ""
}

# create a resource group if it doesn't exist
resource "azurerm_resource_group" "coreosterraform" {
    name = "${var.resourcesname}"
    location = "${var.azurelocation}"
}

# create virtual network
resource "azurerm_virtual_network" "coreosterraformnetwork" {
    name = "tfvn"
    address_space = ["10.0.0.0/16"]
    location = "${var.azurelocation}"
    resource_group_name = "${azurerm_resource_group.coreosterraform.name}"
}

# create subnet
resource "azurerm_subnet" "coreosterraformsubnet" {
    name = "tfsub"
    resource_group_name = "${azurerm_resource_group.coreosterraform.name}"
    virtual_network_name = "${azurerm_virtual_network.coreosterraformnetwork.name}"
    address_prefix = "10.0.2.0/24"
}


# create public IPs
resource "azurerm_public_ip" "coreosterraformips" {
    name = "coreosterraformips"
    location = "${var.azurelocation}"
    resource_group_name = "${azurerm_resource_group.coreosterraform.name}"
    public_ip_address_allocation = "dynamic"

    tags {
        environment = "TerraformDemo"
    }
}

# create network interface
resource "azurerm_network_interface" "coreosterraformnic" {
    name = "tfni"
    location = "${var.azurelocation}"
    resource_group_name = "${azurerm_resource_group.coreosterraform.name}"

    ip_configuration {
        name = "testconfiguration1"
        subnet_id = "${azurerm_subnet.coreosterraformsubnet.id}"
        private_ip_address_allocation = "static"
        private_ip_address = "10.0.2.5"
        public_ip_address_id = "${azurerm_public_ip.coreosterraformips.id}"
    }
}


# create storage account
resource "azurerm_storage_account" "coreosterraformstorage" {
    name                = "jztfcoreosstorage"
    resource_group_name = "${azurerm_resource_group.coreosterraform.name}"
    location = "${var.azurelocation}"
    account_type = "Standard_LRS"

    tags {
        environment = "staging"
    }
}

# create storage container
resource "azurerm_storage_container" "coreosterraformstoragestoragecontainer" {
    name = "vhd"
    resource_group_name = "${azurerm_resource_group.coreosterraform.name}"
    storage_account_name = "${azurerm_storage_account.coreosterraformstorage.name}"
    container_access_type = "private"
    depends_on = ["azurerm_storage_account.coreosterraformstorage"]
}

# create virtual machine
resource "azurerm_virtual_machine" "coreosterraformvm" {
    name = "coreosterraformvm"
    location = "${var.azurelocation}"
    resource_group_name = "${azurerm_resource_group.coreosterraform.name}"
    network_interface_ids = ["${azurerm_network_interface.coreosterraformnic.id}"]
    vm_size = "Standard_A0"

    storage_image_reference {
        publisher = "CoreOS"
        offer = "CoreOS"
        sku = "Stable"
        version = "latest"
    }

    storage_os_disk {
        name = "myosdisk"
        vhd_uri = "${azurerm_storage_account.coreosterraformstorage.primary_blob_endpoint}${azurerm_storage_container.coreosterraformstoragestoragecontainer.name}/myosdisk.vhd"
        caching = "ReadWrite"
        create_option = "FromImage"
    }

    os_profile {
        computer_name = "jzhostname"
        admin_username = "${var.admin_username}"
        admin_password = "JZPassword1234!"
        custom_data = "${base64encode(file("${path.module}/userdata.yml"))}"
    }

    os_profile_linux_config {
        //disable_password_authentication = false
        disable_password_authentication = true
        ssh_keys {
            path = "/home/${var.admin_username}/.ssh/authorized_keys"
            key_data = "ssh_key"
        }
    }

    tags {
        environment = "staging"
    }
}
