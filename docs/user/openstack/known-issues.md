# Known Issues and Workarounds

During this release, we are shipping with a few known issues. We are documenting them here, along with whatever workarounds we are aware of, and have attached the links to where the engineering team is tracking the issues. As changes occur, we will update both this document and the issue trackers with the latest information.

## Long Cluster Names

If the mDNS service name of a server is too long, it will exceed the character limit and cause the installer to fail. To prevent this from happening, please restrict the `metadata.name` field in the `install-config.yaml` to 14 characters. The installer validates this in your install config and throws an error to prevent you from triggering this install time bug. This is being tracked in this [github issue](https://github.com/openshift/installer/issues/2243).

## Resources With Duplicate Names

Since the installer requires the *Name* of your external network and Red Hat Core OS image, if you have other networks or images with the same name, it will choose one randomly from the set. This is not a reliable way to run the installer. We highly recommend that you resolve this with your cluster administrator by creating unique names for your resources in openstack.

## Self Signed Certificates

Due to Terraform not being up to date with Ignition v2.2.0, we are unable to use the installer infrastructure to pass Certificate Authority Bundles to Ignition on Master Nodes. This means that the bootstrap node will be unable to retrieve the ignition configs from swift if your endpoint uses self-signed certificates. As a result, using the AdditionalTrustBundle Field in the install config will not automatically work. What we have observed in testing is that the cert bundles are in fact put in the correct file directories, however, it seems that ignition fails to detect/utilize them. This bug is currently being tracked in [this bugzilla]( https://bugzilla.redhat.com/show_bug.cgi?id=1735192).

## External Network Overlap

If your external network's CIDR range is the same as one of the default network ranges, then you will need to change the matching network range by running the installer with a custom install-config.yaml. If users are experiencing unusual networking problems, please contact your cluster administrator and validate that none of your network CIDRs are overlapping with the external network. We were unfortunately unable to support validation for this due to a lack of support in gophercloud, and even if we were, it is likely that the CIDR range of the floating ip would only be accessable cluster administrators. The default network CIDR
are as follows:

```txt
machineCIDR:    10.0.0.0/16
serviceNetwork: 172.30.0.0/16
clusterNetwork: 10.128.0.0/14
```
