1. Create an install-config.yaml.
The following example uses the settings for the dev cluster.

```
apiVersion: v1beta4
baseDomain: devcluster.openshift.com
metadata:
  name: ${YOUR_OPENSHIFT_CLUSTER_NAME}
networking:
  machineCIDR: "139.178.89.192/26"
platform:
  vsphere:
    virtualCenters:
    - name: vcsa.vmware.devcluster.openshift.com
      username: ${YOUR_VSPHERE_USER}
      password: ${YOUR_VSPHERE_PASSWORD}
      datacenters:
      - dc1
    workspace:
      defaultDatastore: nvme-ds1
    scsiControllerType: pvscsi
    publicNetwork: VM Network
pullSecret: YOUR_PULL_SECRET
sshKey: YOUR_SSH_KEY
```

2. Run `openshift-install create ignition-configs`.

3. Fill out a terraform.tfvars file with the ignition configs generated.
There is an example terraform.tfvars file in this directory named terraform.tfvars.example. The example file is set up for use with the dev cluster running at vcsa.vmware.devcluster.openshift.com. At a minimum, you need to set values for the following variables.
* *cluster_id*. This must match the name used in install-config.yaml.
* *cluster_domain*. This must be ${cluster_id}.${base_domain}.
* *vsphere_user*
* *vsphere_password*
* *ipam_token*
* *bootstrap_ignition_url*. This must be accessible by the bootstrap machine. For example, you could store the bootstrap ignition config in a gist.
* *control_plane_ignition*
* *compute_ignition*

4. Run `terraform init`.

5. Ensure that you have you AWS profile set and a region specified. The installation will use create AWS route53 resources for routing to the OpenShift cluster.

6. Run `terraform apply -auto-approve`.
This will create the OpenShift cluster

7. Run `openshift-install upi bootstrap-complete`. Wait for the bootstrapping to complete.

8. Run `terraform apply -auto-approve -var 'bootstrap_complete=true'`.
This will destroy the bootstrap VM.

9. Run `openshift-install upi finish`. Wait for the cluster install to finish.

10. Enjoy your new OpenShift cluster.

11. Run `terraform destroy -auto-approve`.
