1. Create an install-config.yaml.
There is not a vsphere platform yet, so use the none platform.
The machine CIDR for the dev cluster is 139.178.89.192/26.

```
apiVersion: v1beta4
baseDomain: devcluster.openshift.com
metadata:
  name: mstaeble
networking:
  machineCIDR: "139.178.89.192/26"
platform:
  none: {}
pullSecret: YOUR_PULL_SECRET
sshKey: YOUR_SSH_KEY
```

2. Run `openshift-install create ignition-configs`.

3. Fill out a terraform.tfvars file with the ignition configs generated.
There is an example terraform.tfvars file in this directory named terraform.tfvars.example. The example file is set up for use with the dev cluster running at vcsa.vmware.devcluster.openshift.com. At a minimum, you need to set values for `cluster_id`, `cluster_domain`, `vsphere_user`, `vsphere_password`, `bootstrap_ignition_url`, `control_plane_ignition`, and `compute_ignition`.
The bootstrap ignition config must be placed in a location that will be accessible by the bootstrap machine. For example, you could store the bootstrap ignition config in a gist.

4. Run `terraform init`.

5. Ensure that you have you AWS profile set and a region specified. The installation will use create AWS route53 resources for routing to the OpenShift cluster.

6. Run `terraform apply -auto-approve -var 'step=1'`.
This will create the bootstrap VM.

7. Run `terraform apply -auto-approve -var 'step=2'`.
This will create the control-plane and compute VMs.

8. Run `openshift-install upi bootstrap-complete`. Wait for the bootstrapping to complete.

9. Run `terraform apply -auto-approve -var 'step=3'`.
This will destroy the bootstrap VM.

10. Run `openshift-install upi finish`. Wait for the cluster install to finish.

11. Enjoy your new OpenShift cluster.

12. Run `terraform destroy -auto-approve -var 'step=3'`.
