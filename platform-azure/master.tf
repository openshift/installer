# create virtual network
resource "azurerm_virtual_network" "tectonic_vnet" {
    name = "tectonic_vnet"
    address_space = ["10.0.0.0/16"]
    location = "${var.tectonic_region}"
    resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
}

# create subnet
resource "azurerm_subnet" "tectonic_subnet" {
    name = "tectonic_subnet"
    resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
    virtual_network_name = "${azurerm_virtual_network.tectonic_vnet.name}"
    address_prefix = "10.0.2.0/24"
}


# create public IPs
resource "azurerm_public_ip" "tectonic_ips" {
    name = "tectonic_ips"
    location = "${var.tectonic_region}"
    resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
    public_ip_address_allocation = "dynamic"

    tags {
        environment = "TerraformDemo"
    }
}

# create network interface
resource "azurerm_network_interface" "tectonic_nic" {
    name = "tectonic_nic"
    location = "${var.tectonic_region}"
    resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"

    ip_configuration {
        name = "tectonic_configuration"
        subnet_id = "${azurerm_subnet.tectonic_subnet.id}"
        private_ip_address_allocation = "dynamic"
        load_balancer_backend_address_pools_ids = ["${azurerm_lb_backend_address_pool.k8-lb.id}"]
    }
}


# create storage account
resource "azurerm_storage_account" "tectonic_storage" {
    name                = "jztectonicstorage"
    resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
    location = "${var.tectonic_region}"
    account_type = "Standard_LRS"

    tags {
        environment = "staging"
    }
}

# create storage container
resource "azurerm_storage_container" "tectonic_storage_container" {
    name = "vhd"
    resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
    storage_account_name = "${azurerm_storage_account.tectonic_storage.name}"
    container_access_type = "private"
    depends_on = ["azurerm_storage_account.tectonic_storage"]
}

# create virtual machine
resource "azurerm_virtual_machine" "tectonic_master_vm" {
    name = "tectonic_master_vm"
    location = "${var.tectonic_region}"
    resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
    network_interface_ids = ["${azurerm_network_interface.tectonic_nic.id}"]
    vm_size = "${var.tectonic_azure_vm_size}"

    storage_image_reference {
        publisher = "CoreOS"
        offer = "CoreOS"
        sku = "Stable"
        version = "latest"
    }

    storage_os_disk {
        name = "myosdisk"
        vhd_uri = "${azurerm_storage_account.tectonic_storage.primary_blob_endpoint}${azurerm_storage_container.tectonic_storage_container.name}/myosdisk.vhd"
        caching = "ReadWrite"
        create_option = "FromImage"
    }

    os_profile {
        computer_name = "jzhostname"
        admin_username = "jimzim"
        admin_password = "JZPassword1234!"
        custom_data = "${base64encode("${ignition_config.master.rendered}")}"
    }

    os_profile_linux_config {
        //disable_password_authentication = false
        disable_password_authentication = true
        ssh_keys {
            path = "/home/jimzim/.ssh/authorized_keys"
            key_data = ""
        }
    }

    tags {
        environment = "staging"
    }
}




