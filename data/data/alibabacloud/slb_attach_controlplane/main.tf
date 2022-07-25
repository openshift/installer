provider "alicloud" {
  access_key = var.ali_access_key
  secret_key = var.ali_secret_key
  region     = var.ali_region_id
}

/*
This attachment will be called after the teardown of the bootstrap stage occurs. There is an issue
that occurs during the boostrap teardown that removes all backend servers from the SLB. This adds controlplane
servers into the the list of backend servers for the SLB. The ali_bootstrap_lb defaults to true and when cleanup
is called we set ali_bootstrap_lb to false so that it is removed from the SLB backend.
*/
resource "alicloud_slb_backend_server" "slb_attach_controlplane" {
  count            = length(var.slb_ids)
  load_balancer_id = var.slb_ids[count.index]

  dynamic "backend_servers" {

    for_each = var.ali_bootstrap_lb ? concat([var.bootstrap_ecs_id], var.master_ecs_ids) : var.master_ecs_ids

    content {
      server_id = backend_servers.value
      weight    = 90
    }
  }
}
