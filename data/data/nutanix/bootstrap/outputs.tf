output "bootstrap_ip" {
  value = nutanix_virtual_machine.vm_bootstrap.nic_list_status.0.ip_endpoint_list.0.ip
}
