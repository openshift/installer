# Troubleshooting Tectonic

This directory contains documents about troubleshooting Tectonic clusters.

* [Troubleshooting Tectonic Installer][installer-terraform] describes troubleshooting the installation process itself, including reapplying after a partially completed run, and errors related to missing Terraform modules and failed `tectonic.service` unit.
* [Troubleshooting worker nodes using SSH][worker-nodes] describes how to SSH into a master or worker node to troubleshoot at the host level.
* [Troubleshooting master nodes using SSH][master-nodes] describes how to SSH into a master node to troubleshoot at the host level.
* [Troubleshooting etcd nodes using SSH][etcd-nodes] describes how to SSH into a etcd node to troubleshoot at the host level.
* [Disaster recovery of Scheduler and Controller Manager pods][controller-recovery] describes how to recover a Kubernetes cluster from the failure of certain control plane components.
* [Etcd snapshot troubleshooting][etcd-backup-restore] explains how to spin up a local Kubernetes API server from a backup of another cluster's etcd state for troubleshooting.
* The [Tectonic FAQ][faq] answers some common questions about Tectonic versioning, licensing, and other general matters.


[controller-recovery]: controller-recovery.md
[etcd-backup-restore]: workstation-etcd-api-server-restore.md
[faq]: faq.md
[installer-terraform]: installer-terraform.md
[worker-nodes]: worker-nodes.md
[master-nodes]: master-nodes.md
[etcd-nodes]: etcd-nodes.md
