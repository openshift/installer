# Known Issues and Workarounds

We have been tracking a few issues and FAQs from our users, and are documenting them here along with the known workarounds and solutions. For issues that still have open bugs, we have attached the links to where the engineering team is tracking their progress. As changes occur, we will update both this document and the issue trackers with the latest information.

## Long Cluster Names

If the mDNS service name of a server is too long, it will exceed the character limit and cause the installer to fail. To prevent this from happening, please restrict the `metadata.name` field in the `install-config.yaml` to 14 characters. The installer validates this in your install config and throws an error to prevent you from triggering this install time bug. This is being tracked in this [github issue](https://github.com/openshift/installer/issues/2243).

## Resources With Duplicate Names

Since the installer requires the *Name* of your external network and Red Hat Core OS image, if you have other networks or images with the same name, it will choose one randomly from the set. This is not a reliable way to run the installer. We highly recommend that you resolve this with your cluster administrator by creating unique names for your resources in openstack.

## Extended installation times

Depending on the infrastructure performance, the installation may take longer than what the global installer timeout expects. In those cases, the installer will fail, but the cluster might still converge to a working state. In case of timeout, if such a case is suspected, it is advised to check the cluster health manually after some time:

```shell
$ openshift-install wait-for install-complete
```

## External Network Overlap

If your external network's CIDR range is the same as one of the default network ranges, then you will need to change the matching network range by running the installer with a custom `install-config.yaml`. If users are experiencing unusual networking problems, please contact your cluster administrator and validate that none of your network CIDRs are overlapping with the external network. We were unfortunately unable to support validation for this due to a lack of support in gophercloud, and even if we were, it is likely that the CIDR range of the floating ip would only be accessible cluster administrators. The default network CIDR are as follows:

```txt
machineNetwork:
- cidr: "10.0.0.0/16"
serviceNetwork:
- "172.30.0.0/16"
clusterNetwork:
- cidr: "10.128.0.0/14"
```

## Lack of default DNS servers on created subnets

Some OpenStack clouds do not set default DNS servers for the newly created subnets. In this case, the bootstrap node may fail to resolve public name records to download the OpenShift images or resolve the OpenStack API endpoints.

If you are having this problem in the IPI installer, you will need to set the [`externalDNS` property in `install-config.yaml`](./customization.md#cluster-scoped-properties).

Alternatively, for UPI, you will need to [set the subnet DNS resolvers](./install_upi.md#subnet-dns-optional).

## Deleting machine when instance stuck in provisioning state

The machine controller can get stuck in a delete loop when it tries to delete a machine that is stuck in the provisioning state in OpenStack. This is a bug with OpenStack
because despite the instance existing, it returns a `Not Found` error when the controller attempts to delete it. If you have determined that the proper course of action is to delete the machine, you will first have to manually remove any [finalizers](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#finalizers) from that machine's object. This can be done with the `oc edit` tool. Machines can be found and edited like this:

```sh
oc get machine -n openshift-machine-api
oc edit machine -n openshift-machine-api <name>
```

Once the finalizers are removed, you can delete it.

```sh
oc delete machine -n openshift-machine-api <machine>
```

## Cinder availability zones

OpenShift does not currently support Cinder availability zones. When attaching a volume to a Nova machine, the Cloud Provider will look for available storage in the same Availability Zone (or better said: in a Cinder availability Zone with the same name as the Nova availability zone of the corresponding machine).

In 4.6, [it is possible to control what Availability Zone each machine will be created in][nova-az-setting]. Cloud Provider can be instructed to ignore the corresponding machine's AZ (and thus pick storage regardless of the zones) by adding the `ignore-volume-az = yes` directive in its configuration, under the `[BlockStorage]` section of the `cloud-provider-config` configmap:

```sh
oc edit cm cloud-provider-config -n openshift-config
```

```txt
[BlockStorage]
ignore-volume-az = yes
```

Then it is required to create and use a new Storage Class that has the `availability` parameter set to your Cinder availability zone name, as below:

```txt
allowVolumeExpansion: true
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: sc-availability
provisioner: kubernetes.io/cinder
parameters:
  availability: <Cinder_AZ>
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
```

[nova-az-setting]: ../openstack#setting-nova-availability-zones

## Problems when changing the Machine spec for master node

Changing the Machine spec (e.g. the image name) for master node will delete the instance in OpenStack and the machine status will change to `phase: failed with errorMessage: Can't find created instance`.
This is a known limitation in cluster-api-provider-openstack where we don't support that action.
To recover from a scenario where one or multiple master nodes are in failed status, you will have to follow the [restore procedure](https://docs.openshift.com/container-platform/4.6/backup_and_restore/replacing-unhealthy-etcd-member.html#restore-replace-stopped-etcd-member_replacing-unhealthy-etcd-member) to scale-up master nodes to reach quorum again.

Even if master node was removed in Nova, it can still appear when you run `oc get machines -n openshift-machine-api`.
To delete it, you can run: `oc delete machine --force=true -n openshift-machine-api <$NODE-NAME>`

# Known Issues specific to User-Provisioned Installations

## Stale resources

The teardown playbooks provided for UPI installation will not delete:
 - Cinder volumes from PVs
 - Swift container for image registry (bootstrap container is correctly deleted)

These objects have to be manually removed after running the teardown playbooks.

## Control plane scale out anti-affinity

A long-standing OpenStack Compute issue prevents the "soft"-anti-affinity
policy to be correctly applied when instances are created in parallel. The
Installer places the initial control plane nodes in a Server group with
"soft-anti-affinity" policy and creates the first three sequentially to work
around the problem. Additional control plane nodes beyond the first three are
created, according to the "replicas" property, in a later phase. These nodes
are not guaranteed to be created sequentially.

When anti-affinity is a requirement, it is advised to rely on the more stable
"anti-affinity" for any node beyond the third. To do so, install with three
control plane replicas only. Add more nodes as a day-2 operation by assigning
them a new Server group with "anti-affinity" policy. Use the `ServerGroupID`
property of the Machine ProviderSpec.

## Requirement to create Control Plane Machines manifests (Kuryr SDN)

Installations with Kuryr SDN can timeout due to changes in the way Kuryr detects
the OpenStack Subnet used by the cluster's nodes. Kuryr relied on the Network of
the cluster's nodes Subnet having a specific tag, but the tag was removed for IPI
Installations causing the need to discover it from the OpenShift Machine objects,
which the creation is removed on one of the UPI steps. Until the fix for
[the issue][bugzilla-upi] is available, as a workaround, only the compute machine
manifests should be removed in the [Remove machines and machinesets][manifests-removal]
section of the UPI guide. The command to run is:

```console
$ rm -f openshift/99_openshift-cluster-api_worker-machineset-*.yaml
```
[bugzilla-upi]: https://bugzilla.redhat.com/show_bug.cgi?id=1927244
[manifests-removal]:../openstack/install_upi.md#remove-machines-and-machinesets
