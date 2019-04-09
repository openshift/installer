1. Create an install-config.yaml.
The following example uses the settings for the dev cluster.

```
apiVersion: v1beta4
baseDomain: devcluster.openshift.com
metadata:
  name: YOUR_CLUSTER_NAME
networking:
  machineCIDR: "139.178.89.192/26"
platform:
  vsphere:
    vCenter: vcsa.vmware.devcluster.openshift.com
    username: YOUR_VSPHERE_USER
    password: YOUR_VSPHERE_PASSWORD
    datacenter: dc1
    defaultDatastore: nvme-ds1
pullSecret: YOUR_PULL_SECRET
sshKey: YOUR_SSH_KEY
```

2. Run `openshift-install create ignition-configs`.

3. Run `IPAM_TOKEN=YOUR_TOKEN BOOTSTRAP_IGNITION_URL=YOUR_URL create_tfvars.sh`
This needs to be run in your asset directory so that the results from the installer can be used to fill out terraform.tfvars.
The bootstrap ignition config must be placed in a location that will be accessible by the bootstrap machine. For example, you could store the bootstrap ignition config in a gist.

4. Run `terraform init`.

5. Ensure that you have you AWS profile set and a region specified. The installation will use create AWS route53 resources for routing to the OpenShift cluster.

6. Run `terraform apply -auto-approve`.
This will reserve IP addresses for the VMs.

7. Run `openshift-install wait-for bootstrap-complete`. Wait for the bootstrapping to complete.

8. Run `terraform apply -auto-approve -var 'bootstrap_complete=true'`.
This will destroy the bootstrap VM.

9. Run `openshift-install wait-for install-complete`. Wait for the cluster install to finish.

10. Enjoy your new OpenShift cluster.

11. Run `terraform destroy -auto-approve`.
