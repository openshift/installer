
locals {
  vhd-converter_nic_ip_configuration_name = "vhd-converter-nic-ip"
}

resource "random_string" "username" {
  length  = 16
  upper   = true
  special = false
  number  = true
}

resource "random_string" "password" {
  length  = 32
  upper   = true
  special = true
  override_special = "!_-"
  number  = true
}

data "azurerm_storage_account_sas" "vhd-file" {
  connection_string = var.storage_account.primary_connection_string
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = timestamp()
  expiry = timeadd(timestamp(), "3h")

  permissions {
    read    = true
    list    = false
    create  = true
    add     = false
    delete  = false
    process = false
    write   = true
    update  = true
  }
}

resource "azurerm_public_ip" "vhd-converter_public_ip" {
  sku                 = "Standard"
  location            = var.azure_region
  name                = "${var.cluster_id}-vhd-converter-pip"
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "vhd-converter" {
  name                = "${var.cluster_id}-vhd-converter-nic"
  location            = var.azure_region
  resource_group_name = var.resource_group_name

  ip_configuration {
    subnet_id                     = var.subnet_id
    name                          = local.vhd-converter_nic_ip_configuration_name
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vhd-converter_public_ip.id
  }
}

# Rendere the bash script template with Terraform variables.
data "template_file" "process-template" {
  template = file("${path.module}/convert.sh.tpl")
  vars = {
    sas_token = data.azurerm_storage_account_sas.vhd-file.sas
    primary_blob_endpoint = var.storage_account.primary_blob_endpoint
    container_name = var.container_name
    vhd_url = var.azure_image_url
  }
}

resource "local_file" "bash-script" {
  filename = "${path.module}/convert.sh"
  content = data.template_file.process-template.rendered
}

resource "azurerm_virtual_machine" "vhd-converter" {
  name                  = "${var.cluster_id}-vhd-converter"
  location              = var.azure_region
  resource_group_name   = var.resource_group_name
  network_interface_ids = [azurerm_network_interface.vhd-converter.id]
  vm_size               = var.vm_size

  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  identity {
    type         = "UserAssigned"
    identity_ids = [var.identity]
  }

  storage_os_disk {
    name              = "${var.cluster_id}-vhd-converter_OSDisk" # os disk name needs to match cluster-api convention
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
    disk_size_gb      = 100
  }

  # Some linux image from the Azure marketplace.
  storage_image_reference {
      publisher = "canonical"
      offer     = "UbuntuServer"      
      sku = "18.04-LTS"
      version = "latest"
  }

  os_profile {
    computer_name  = "${var.cluster_id}-vhd-converter-vm"
    admin_username = random_string.username.result
    admin_password = random_string.password.result
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  boot_diagnostics {
    enabled     = false
    storage_uri = var.storage_account.primary_blob_endpoint
  }

  # Copy VHD conversion script in VM.
  provisioner "file" {
    connection {
      type = "ssh"
      host = azurerm_public_ip.vhd-converter_public_ip.ip_address
      agent = false
      user = random_string.username.result
      password = random_string.password.result
      timeout = "180s"
    }

    source = "${path.module}/convert.sh"
    destination = "/home/${random_string.username.result}/convert.sh"
  }

  # Start VHD conversion and wait until it finishes. After the script ran successfully, the VM is marked as completed
  # and Terraform will start creating the VM image from the VHD file uploaded to the blob storage.
  # Without 'remote-exec' Terraform wouldn't wait before trying to create the VM image. It would fail.
  provisioner "remote-exec" {
    connection {
      type = "ssh"
      host = azurerm_public_ip.vhd-converter_public_ip.ip_address
      agent = false
      user = random_string.username.result
      password = random_string.password.result
      timeout = "180s"
    }

    inline = [
      "sudo chmod +x ./convert.sh",
      "./convert.sh",
    ]
  }  
}

resource "azurerm_network_security_rule" "vhd-converter_ssh_in" {
  name                        = "vhd-converter_ssh_in"
  priority                    = 105
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = var.resource_group_name
  network_security_group_name = var.nsg_name
}
