1. Create an install-config.yaml.
The machine CIDR for the dev cluster is 139.178.89.192/26.

```
apiVersion: v1beta4
baseDomain: devcluster.openshift.com
metadata:
  name: mstaeble
networking:
  machineCIDR: "139.178.89.192/26"
platform:
  vsphere: {}
pullSecret: YOUR_PULL_SECRET
sshKey: YOUR_SSH_KEY
```

2. Run `openshift-install create ignition-configs`.

3. Fill out a terraform.tfvars file with the ignition configs generated.
There is an example terraform.tfvars file in this directory named terraform.tfvars.example. The example file is set up for use with the dev cluster running at vcsa.vmware.devcluster.openshift.com. At a minimum, you need to set values for `cluster_id`, `cluster_domain`, `vsphere_user`, `vsphere_password`, `pull_secret`, `bootstrap_ignition_url`, `control_plane_ignition`, and `compute_ignition`.
The bootstrap ignition config must be placed in a location that will be accessible by the bootstrap machine. For example, you could store the bootstrap ignition config in a gist.
Initially, the `bootstrap_complete` variable must be false, the `bootstrap_ip` variable must be an empty string, and the `control_plane_ips variable must be an empty list.
To secure your pull secret, you should remove the pull from the bootstrap ignition config and pass it as a variable to terraform.
  a) Create an ignition config without the pull secret via `jq 'del(.storage.files[] | select(.path=="/root/.docker/config.json"))' bootstrap.ign`.
  b) Extract the pull secret to pass to terraform via `jq '.storage.files[] | select(.path=="/root/.docker/config.json")' bootstrap.ign`.

4. Run `terraform init`.

5. Run `terraform apply -auto-approve`.

6. Find the IP address of the bootstrap machine.
If you provided an extra user, you can use that user to log into the bootstrap machine via the vSphere web console.
Alternatively, you could iterate through the IP addresses in the 139.178.89.192/26 block looking for one that has the expected hostname, which is bootstrap-0.{cluster_domain}. For example, `ssh -i ~/.ssh/libra.pem -o StrictHostNameChecking=no -q core@139.178.89.199 hostname`

7. Update the terraform.tfvars file with the IP address of the bootstrap machine.

8. Run `terraform apply -auto-approve`.
From this point forward, route53 resources will be managed by terraform. You will need to have your AWS profile set and a region specified.

9. Find the IP addresses of the control plane machines. See step 6 for examples of how to do this. The expected hostnames are control-plane-{0,1,2}.{cluster_domain}. The control plane machines will change their IP addresses once. You need the final IP addresses. If you happen to use the first set of IP addresses, you can later update the IP addresses in the terraform.tfvars file and re-run terraform.

10. Update the terraform.tfvars file with the IP addresses of the control plane machines.

11. Run `terraform apply -auto-approve`.

12. Run `openshift-install user-provided-infrastructure`. Wait for the bootstrapping to complete.
You *may* need to log into each of the control plane machines. It would seem that, for some reason, the etcd-member pod does not start until the machine is logged into.

13. Update the terraform.tfvars file to set the `bootstrap_complete` variable to "true".

14. Run `terraform apply -auto-approve`.

15. Run `openshift-install user-provided-infrastructure finish`. Wait for the cluster install to finish.
Currently, the cluster install does not finish. There is an outstanding issue with the openshift-console operator not installing successfully. The cluster should still be usable save for the console, however.

16. Enjoy your new OpenShift cluster.

17. Run `terraform destroy -auto-approve`.
