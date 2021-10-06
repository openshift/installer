output "master_ips" {
  value = data.ibm_pi_instance_ip.master_ip.*.ip
}
