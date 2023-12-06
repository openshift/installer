# Known Issues and Workarounds

We have been tracking a few issues and FAQs from our users, and are documenting them here along with the known workarounds and solutions. For issues that still have open bugs, we have attached the links to where the engineering team is tracking their progress. As changes occur, we will update both this document and the issue trackers with the latest information.

## HTTPS CommonName deprecation

With OpenShift v4.10, HTTPS certificates must include the names for the server in the `Subject Alternative Names` field. Legacy certificates will be rejected with the following error message:

> x509: certificate relies on legacy Common Name field, use SANs instead

In order to validate your OpenStack infrastructure prior to installation or upgrade to v4.10, please refer to the [dedicated instructions](invalid-https-certificates.md).

## Resources With Duplicate Names

Since the installer requires the *Name* of your external network and Red Hat Core OS image, if you have other networks or images with the same name, it will choose one randomly from the set. This is not a reliable way to run the installer. We highly recommend that you resolve this with your cluster administrator by creating unique names for your resources in openstack.

## Soft-anti-affinity

A long-standing OpenStack Compute issue prevents the soft anti-affinity
policy to correctly apply when instances are created in parallel. To maximize
the chances that the default soft-anti-affinity setting applies, the Installer
generates the Control plane instances sequentially. However, there is no
guarantee that other machines, and in particular Compute nodes, are created
sequentially. This applies both to install-time machine building, and day 2
operations (e.g. MachineSet scale-out).

**When it is a requirement that machines land on distinct Compute hosts,
explicitly set `serverGroupPolicy` to `anti-affinity`.**

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

## Changing user domain for image registry

When users change their user domain name or ID globally they hit an issue with the image registry, when the operator doesn't update its configuration accordingly. It leads to the fact that the operator can't authenticate in OpenStack Swift and becomes degraded.

```sh
$ oc logs cluster-image-registry-operator-748f8f9855-fhblt -c cluster-image-registry-operator --tail=20

I0303 17:32:15.100153      14 controller.go:291] object changed: *v1.Config, Name=cluster (status=true): changed:status.conditions.0.lastTransitionTime={"2021-03-03T17:32:14Z" -> "2021-03-03T17:32:15Z"}
I0303 17:32:15.603731      14 controller.go:291] object changed: *v1.Config, Name=cluster (status=true): 
E0303 17:32:15.610609      14 controller.go:330] unable to sync: unable to sync storage configuration: Failed to authenticate provider client: Authentication failed, requeuing
I0303 17:32:15.644501      14 controller.go:291] object changed: *v1.Config, Name=cluster (status=true): 
E0303 17:32:15.650130      14 controller.go:330] unable to sync: unable to sync storage configuration: Failed to authenticate provider client: Authentication failed, requeuing
```

This issue was fixed in 4.8, which means that all new deployments work correctly and update image registry configuration in the case of credentials change. For older versions, and for clusters that have been upgraded to 4.8, you need to manually clean the config and delete the incorrect values. 

To do so, execute `oc edit configs.imageregistry.operator.openshift.io/cluster` and manually remove `domain` or `domainID` in the `.spec.storage.swift` section. Alternatively you can use the `oc patch` command:

```sh
$ oc patch configs.imageregistry.operator.openshift.io/cluster --type 'json' -p='[{"op": "remove", "path": "/spec/storage/swift/domain"}]'
config.imageregistry.operator.openshift.io/cluster patched
```

# Known Issues specific to User-Provisioned Installations

## Stale resources

The teardown playbooks provided for UPI installation will not delete:
 - Cinder volumes from PVs
 - Swift container for image registry (bootstrap container is correctly deleted)

These objects have to be manually removed after running the teardown playbooks.

## Limitations of creating external load balancers using pre-defined FIPs

On most clouds, the default policy prevents non-admin users from creating                                     
a floating IP that has a specific address. Such policies cause the cloud provider
to fail handling floating IP assignment to load balancers if a floating IP address
is present in the service specification.

You can pre-create a floating IP and pass the address of it in the service
specification if using the external cloud provider. The in-tree cloud
provider does not support this.
Alternatively, you can [relax the Neutron policy][policy-change-steps]
to allow non-admin users to create FIPs with a specific IP address.

[policy-change-steps]: https://access.redhat.com/solutions/6069071
