output "master-vm-ids" {
  value = ["${azurerm_virtual_machine.tectonic_master.*.id}"]
}
