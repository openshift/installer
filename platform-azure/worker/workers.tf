# TODO
# Add to availabilityset

# create network interface
resource "azurerm_network_interface" "tectonic_agent_nic" {
    name = "tectonic_agent_nic-${count.index}"
    location = "${var.tectonic_azure_location}"
    resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"

    ip_configuration {
        name = "tectonic_configuration"
        subnet_id = "${azurerm_subnet.tectonic_subnet.id}"
        private_ip_address_allocation = "dynamic"
        load_balancer_backend_address_pools_ids = ["${azurerm_lb_backend_address_pool.k8-lb.id}"]
    }
}

resource "azurerm_virtual_machine" "tectonic_worker_vm" {
    name = "tectonic_worker_vm-${count.index}"
    location = "${var.tectonic_azure_location}"
    resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
    network_interface_ids = ["${azurerm_network_interface.tectonic_agent_nic.id}"]
    vm_size = "${var.tectonic_azure_vm_size}"

    storage_image_reference {
        publisher = "CoreOS"
        offer = "CoreOS"
        sku = "Stable"
        version = "latest"
    }

    storage_os_disk {
        name = "${azurerm_virtual_machine.tectonic_worker_vm.name}-osdisk"
        vhd_uri = "${azurerm_storage_account.tectonic_storage.primary_blob_endpoint}${azurerm_storage_container.tectonic_storage_container.name}/${azurerm_virtual_machine.tectonic_worker_vm.name}-osdisk.vhd"
        caching = "ReadWrite"
        create_option = "FromImage"
    }

    os_profile {
        computer_name = "tectonic-worker-${count.index}"
        admin_username = "core"
        admin_password = ""
        custom_data = "${base64encode("${ignition_config.worker.*.rendered[count.index]}")}"
    }

    os_profile_linux_config {
        disable_password_authentication = true
        ssh_keys {
            path = "/home/core/.ssh/authorized_keys"
            key_data = "${file(var.tectonic_ssh_key)}"
        }
    }

    tags {
        environment = "staging"
    }
}
