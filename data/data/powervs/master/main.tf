# TODO(mjturek): network and image data blocks can be in main module
#                as master and bootstrap will be using the same
#                network and image. Once we add in master module, make
#                the move.
data "ibm_pi_network" "network" {
  pi_network_name      = var.network_name
  pi_cloud_instance_id = var.cloud_instance_id
}

data "ibm_pi_image" "master_image" {
  pi_image_name        = var.image_name
  pi_cloud_instance_id = var.cloud_instance_id
}

# Create the master instances
resource "ibm_pi_instance" "master" {
  count                = var.instance_count
  pi_memory            = var.memory
  pi_processors        = var.processors
  pi_instance_name     = "${var.cluster_id}-master-${count.index}"
  pi_proc_type         = var.proc_type
  pi_image_id          = data.ibm_pi_image.master_image.id
  pi_sys_type          = var.sys_type
  pi_cloud_instance_id = var.cloud_instance_id
  pi_network_ids       = [data.ibm_pi_network.network.id]

  pi_user_data     = base64encode(var.ignition)
  pi_key_pair_name = var.key_id
  pi_health_status = "WARNING"
}

data "ibm_pi_instance_ip" "master_ip" {
  count      = var.instance_count
  depends_on = [ibm_pi_instance.master]

  pi_instance_name     = ibm_pi_instance.master[count.index].pi_instance_name
  pi_network_name      = data.ibm_pi_network.network.pi_network_name
  pi_cloud_instance_id = var.cloud_instance_id
}
