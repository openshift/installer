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

# Known Issues specific to User-Provisioned Installations

## Stale resources

The teardown playbooks provided for UPI installation will not delete:
 - Cinder volumes from PVs
 - Swift container for image registry (bootstrap container is correctly deleted)

These objects have to be manually removed after running the teardown playbooks.
