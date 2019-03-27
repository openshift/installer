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
There is an example terraform.tfvars file in this directory named terraform.tfvars.example. The example file is set up for use with the dev cluster running at vcsa.vmware.devcluster.openshift.com. At a minimum, you need to set values for the following variables.
* cluster_id
* cluster_domain
* vsphere_user
* vsphere_password
* bootstrap_ip
* control_plane_ips
* compute_ips
* pull_secret
* bootstrap_ignition_url
* control_plane_ignition
* compute_ignition
The bootstrap ignition config must be placed in a location that will be accessible by the bootstrap machine. For example, you could store the bootstrap ignition config in a gist.
For the IP addresses, you should have static IP addresses reserved for you.
To secure your pull secret, you should remove the pull from the bootstrap ignition config and pass it as a variable to terraform.
  a) Create an ignition config without the pull secret via `jq 'del(.storage.files[] | select(.path=="/root/.docker/config.json"))' bootstrap.ign`.
  b) Extract the pull secret to pass to terraform via `jq '.storage.files[] | select(.path=="/root/.docker/config.json")' bootstrap.ign`.

4. Run `terraform init`.

5. Ensure that you have you AWS profile set and a region specified. The installation will use create AWS route53 resources for routing to the OpenShift cluster.

6. Run `terraform apply -auto-approve`.
This will create the OpenShift cluster

7. Run `openshift-install upi bootstrap-complete`. Wait for the bootstrapping to complete.

8. Run `terraform apply -auto-approve -var 'bootstrap_ip='`.
This will destroy the bootstrap VM.

9. Run `openshift-install upi finish`. Wait for the cluster install to finish.

10. Enjoy your new OpenShift cluster.

11. Run `terraform destroy -auto-approve`.
