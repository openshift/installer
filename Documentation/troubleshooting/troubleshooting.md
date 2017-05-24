# Troubleshooting Tectonic

This folder contains documents specific to troubleshooting Tectonic Installer.

* [Troubleshooting Tectonic Installer][installer-terraform] describes troubleshooting the installation process itself, including reapplying after a partially completed run, and errors related to missing Terraform modules and failed `tectonic.service` unit.
* [Disaster recovery of Scheduler and Controller Manager pods][controller-recovery] describes how to recover a Kubernetes cluster from the failure of certain control plane components.
* [Troubleshooting worker nodes using SSH][worker-nodes] describes how to SSH into a controller or worker node to troubleshoot at the host level.


[installer-terraform]: installer-terraform.md
[controller-recovery]: controller-recovery.md
[worker-nodes]: worker-nodes.md
