# OpenShift 4 installer on OpenStack troubleshooting

Unfortunately, there will always be some cases where OpenShift fails to install properly. In these events, it is helpful to understand the likely failure modes as well as how to troubleshoot the failure.

This document discusses some troubleshooting options for OpenStack based
deployments. For general tips on troubleshooting the installer, see the [Installer Troubleshooting](../troubleshooting.md) guide.

## View instances logs

OpenStack CLI tools should be installed, then:

`openstack console log show <instance>`

## SSH access to the instances

Get the IP address of the node on the private network:
```
openstack server list | grep master
| 0dcd756b-ad80-42f1-987a-1451b1ae95ba | cluster-wbzrr-master-1     | ACTIVE    | cluster-wbzrr-openshift=172.24.0.21                | rhcos           | m1.s2.xlarge |
| 3b455e43-729b-4e64-b3bd-1d4da9996f27 | cluster-wbzrr-master-2     | ACTIVE    | cluster-wbzrr-openshift=172.24.0.18                | rhcos           | m1.s2.xlarge |
| 775898c3-ecc2-41a4-b98b-a4cd5ae56fd0 | cluster-wbzrr-master-0     | ACTIVE    | cluster-wbzrr-openshift=172.24.0.12                | rhcos           | m1.s2.xlarge |
```

And connect to it using the master currently holding the API VIP (and hence the API FIP) as a jumpbox:

```
ssh -J core@${FIP} core@<host>
```
