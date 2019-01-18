# OpenShift 4 installer on OpenStack troubleshooting

Support for launching clusters on OpenStack is **experimental**.

Unfortunately, there will always be some cases where OpenShift fails to install properly. In these events, it is helpful to understand the likely failure modes as well as how to troubleshoot the failure.

This document discusses some troubleshooting options for OpenStack based
deployments. For general tips on troubleshooting the installer, see the [Installer Troubleshooting](../troubleshooting.md) guide.

## View instances logs

OpenStack CLI tools should be installed, then:

`openstack console log show <instance>`

## ssh access to the instances

By default, the only exposed instance is the service VM, but ssh access is not
allowed (nor the ssh key injected), so in case you want to access the hosts, it
is required to create a floating IP and attach it to some host (master-0 in
this example)

### Create security group to allow ssh access

```
INSTANCE=$(openstack server list -f value -c Name | grep master-0)
openstack security group create ssh
# Note this opens port 22/tcp to 0.0.0.0/0
openstack security group rule create \
  --ingress \
  --protocol tcp \
  --dst-port 22 \
  ssh
openstack server add security group ${INSTANCE} ssh
```

Optionally, allow ICMP traffic (to ping the instance):

```
openstack security group rule create \
  --ingress \
  --protocol icmp \
  ssh
```

### Create and attach the floating IP

```
# This must be set to the external network configured in the OpenShift install
PUBLIC_NETWORK="external_network"

INSTANCE=$(openstack server list -f value -c Name | grep master-0)

FIP=$(openstack floating ip create ${PUBLIC_NETWORK} --description ${INSTANCE} -f value -c floating_ip_address)

openstack server add floating ip ${INSTANCE} ${FIP}
```

### Access the host

```
ssh core@${FIP}
```

You can use it as jump host as well:

```
ssh -J core@${FIP} core@<host>
```

NOTE: If you are running the `openshift-installer` from an all-in-one OpenStack
deployment (compute + controller in a single host), you can connect to the
instance network namespace directly:

```
NODE_ADDRESSES=$(openstack server show ${INSTANCE} -f value -c addresses | cut -d',' -f1)
NODE_IP=${NODE_ADDRESSES#"openshift="}
sudo ip netns exec "qdhcp-$(openstack network show openshift -f value -c id)" ssh core@$NODE_IP
```
