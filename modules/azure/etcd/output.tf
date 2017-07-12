output "node_names" {
  value = ["${compact(split(";", var.base_domain == "" ?
    join(";", slice(formatlist("${var.cluster_name}-%s", var.const_internal_node_names), 0, var.etcd_count)) : 
    join(";", formatlist("%s.${var.base_domain}", slice(formatlist("${var.cluster_name}-%s", var.const_internal_node_names), 0, var.etcd_count)))))
  }"]
}
