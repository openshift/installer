locals {
  description = "Created By OpenShift Installer"
  prefix      = var.cluster_id
}

resource "alicloud_instance" "master" {
  count             = length(var.vswitch_ids)
  resource_group_id = var.resource_group_id

  instance_name              = "${local.prefix}_master_${count.index}"
  instance_type              = var.instance_type
  image_id                   = var.image_id
  internet_max_bandwidth_out = 0

  vswitch_id      = var.vswitch_ids[count.index]
  security_groups = [var.sg_id]
  role_name       = var.role_name

  system_disk_name        = "${local.prefix}_sys_disk-master_${count.index}"
  system_disk_description = local.description
  system_disk_category    = var.system_disk_category
  system_disk_size        = var.system_disk_size

  user_data = var.user_data_ign
  key_name  = var.key_name
  tags = merge(
    {
      "Name" = "${local.prefix}-master-${count.index}"
    },
    var.tags,
  )
}

resource "alicloud_slb_attachment" "slb_attachment_master" {
  load_balancer_id = var.slb_id
  instance_ids     = alicloud_instance.master.*.id
  weight           = 90
}
