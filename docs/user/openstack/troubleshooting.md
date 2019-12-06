# OpenShift 4 installer on OpenStack troubleshooting

Unfortunately, there will always be some cases where OpenShift fails to install properly. In these events, it is helpful to understand the likely failure modes as well as how to troubleshoot the failure.

This document discusses some troubleshooting options for OpenStack based
deployments. For general tips on troubleshooting the installer, see the [Installer Troubleshooting](../troubleshooting.md) guide.

## View instances logs

OpenStack CLI tools should be installed, then:

`openstack console log show <instance>`

## Cluster destroying fails

Destroying the cluster has been noticed to [sometimes fail](https://github.com/openshift/installer/issues/1985). We are working on patching this, but in a mean time the workaround is to simply restart the destroying process of the cluster.

## Machine has ERROR state

This could be because the machine's instance was accidentally destroyed and the cluster API provider cannot recreate it.

You can check the status of machines with the help of the command

```sh
oc get machines -n openshift-machine-api
```

If the broken machine is a master then follow the instructions in the [disaster recovery documentation](https://docs.openshift.com/container-platform/4.1/disaster_recovery/scenario-1-infra-recovery.html).

For workers, you should delete the machine manually with

```sh
oc delete machine -n openshift-machine-api <machine_name>
```

The operation can take up to 5 minutes, during which time the machine will be gracefully removed and all its resources returned to the pool.

A new worker machine for the cluster will soon be created automatically by the [machine-api-operator](https://github.com/openshift/machine-api-operator).

**NOTE**: in future versions of OpenShift all broken machines will be automatically deleted and recovered by the machine-api-operator.

## SSH access to the instances

Get the IP address of the node on the private network:

```sh
openstack server list | grep master
| 0dcd756b-ad80-42f1-987a-1451b1ae95ba | cluster-wbzrr-master-1     | ACTIVE    | cluster-wbzrr-openshift=172.24.0.21                | rhcos           | m1.s2.xlarge |
| 3b455e43-729b-4e64-b3bd-1d4da9996f27 | cluster-wbzrr-master-2     | ACTIVE    | cluster-wbzrr-openshift=172.24.0.18                | rhcos           | m1.s2.xlarge |
| 775898c3-ecc2-41a4-b98b-a4cd5ae56fd0 | cluster-wbzrr-master-0     | ACTIVE    | cluster-wbzrr-openshift=172.24.0.12                | rhcos           | m1.s2.xlarge |
```

And connect to it using the master currently holding the API VIP (and hence the API FIP) as a jumpbox:

```sh
ssh -J core@${FIP} core@<host>
```
