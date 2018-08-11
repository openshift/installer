# Etcd
output "etcd_sg_id" {
  value = "${module.sg.etcd_sg_id}"
}

# Masters
output "master_sg_id" {
  value = "${module.sg.master_sg_id}"
}

# Workers
output "worker_sg_id" {
  value = "${module.sg.worker_sg_id}"
}
