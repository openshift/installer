data "ibm_pi_network" "network" {
    pi_network_name         = var.network_name
    pi_cloud_instance_id    = var.cloud_instance_id
}

data "ibm_pi_image" "bootstrap_image" {
    pi_image_name           = var.image_name
    pi_cloud_instance_id    = var.cloud_instance_id
}

resource "ibm_pi_instance" "bootstrap" {
    pi_memory               = var.bootstrap.memory
    pi_processors           = var.bootstrap.processors
    pi_instance_name        = "${var.cluster_id}-bootstrap"
    pi_proc_type            = var.proc_type
    pi_image_id             = data.ibm_pi_image.bootstrap_image.id
    pi_sys_type             = var.sys_type
    pi_cloud_instance_id    = var.cloud_instance_id
    pi_network_ids          = [data.ibm_pi_network.network.id]

    # Not needed by RHCOS but required by resource
    pi_key_pair_name        = "${var.cluster_id}-keypair"
    pi_health_status        = "WARNING"
}
