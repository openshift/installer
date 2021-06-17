locals {
  description = "Created By OpenShift Installer"
  prefix      = var.cluster_id
}

data "alicloud_instances" "master_data" {
  ids = alicloud_instance.master.*.id
}

resource "alicloud_instance" "master" {
  count             = length(var.vswitch_ids)
  resource_group_id = var.resource_group_id

  host_name                  = "${local.prefix}-master${count.index}"
  instance_name              = "${local.prefix}-master${count.index}"
  instance_type              = var.instance_type
  image_id                   = var.image_id
  internet_max_bandwidth_out = 0

  vswitch_id      = var.vswitch_ids[count.index]
  security_groups = [var.sg_id]
  role_name       = var.role_name

  system_disk_name        = "${local.prefix}_sys_disk-master${count.index}"
  system_disk_description = local.description
  system_disk_category    = var.system_disk_category
  system_disk_size        = var.system_disk_size

  user_data = base64encode(var.user_data_ign)
  tags = merge(
    {
      "Name" = "${local.prefix}-master${count.index}"
    },
    var.tags,
  )
}

resource "alicloud_slb_backend_server" "slb_attachment_masters" {
  count            = "${length(var.slb_ids) * length(alicloud_instance.master.*.id)}"
  load_balancer_id = "${element(var.slb_ids, ceil(count.index / length(alicloud_instance.master.*.id)))}"
  backend_servers {
    server_id = "${element(alicloud_instance.master.*.id, count.index)}"
    weight    = 90
  }
}
