output "control_plane_interfaces" {
  value = data.ironic_introspection.openshift-master-introspection.*.interfaces
}
