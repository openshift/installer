output "bootstrap_ip" {
  value = module.vm.bootstrap_ip
}

output "api_member_ext_id" {
  value = module.lb.api_member_ext_id
}

output "api_member_int_id" {
  value = module.lb.api_member_int_id
}
