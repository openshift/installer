output "control_plane_ips" {
  value = <<EOT
    var.master_count > 2 ? [
      data.ironic_introspection.openshift-master-introspection.interfaces.*.ip ] :
        var.master_count > 1 ? [
          data.ironic_introspection.openshift-master-introspection.interfaces.0.ip,
          data.ironic_introspection.openshift-master-introspection.interfaces.1.ip ] :
            var.master_count > 0 ? [
              data.ironic_introspection.openshift-master-introspection.interfaces.0.ip ] : 0
EOT
}
