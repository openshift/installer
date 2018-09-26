# XXX(crawford): This is only needed because the installer will only run either
#                AWS or libvirt. This prevents us from outputting anything
#                directly from the base.
output "ignition_bootstrap" {
  value = "${module.assets_base.ignition_bootstrap}"
}
