# OpenStack Platform Support

Support for launching clusters on OpenStack is **experimental**.

This document discusses the requirements, current expected behavior, and how to
try out what exists so far.

## OpenStack Requirements

The installer assumes the following about the OpenStack cloud you run against:

* You must create a `clouds.yaml` file with the auth URL and credentials
    necessary to access the OpenStack cloud you want to use.  Information on
    this file can be found at
    https://docs.openstack.org/os-client-config/latest/user/configuration.html
    and it looks like:
```
clouds:
  shiftstack:
    auth:
      auth_url: http://10.10.14.42:5000/v3
      project_name: shiftstack
      username: shiftstack_user
      password: XXX
      user_domain_name: Default
      project_domain_name: Default
  dev-evn:
    region_name: RegionOne
    auth:
      username: 'devuser'
      password: XXX
      project_name: 'devonly'
      auth_url: 'https://10.10.14.22:5001/v2.0'
```

* Swift must be enabled.  The user must have `swiftoperator` permissions and
  `temp-url` support must be enabled. As an OpenStack admin:
  * `openstack role add --user <user> --project <project> swiftoperator`
  * `openstack object store account set --property Temp-URL-Key=superkey`

**NOTE:** Swift is required as the user-data provided by OpenStack is not big
enough to store the ignition config files, so they are served by swift instead.

* You may need to increase the security group related quotas from their default
  values. For example (as an OpenStack admin) `openstack quota set --secgroups 100 --secgroup-rules 1000 <project>`

* The installer requires a proper RHCOS image in the OpenStack cluster or project:
`openstack image create --container-format=bare --disk-format=qcow2 --file redhat-coreos-${RHCOSVERSION}-openstack.qcow2 redhat-coreos-${RHCOSVERSION}`

**NOTE:** Depending on your OpenStack environment you can upload the RHCOS image
as `raw` or `qcow2`. See [Disk and container formats for images](https://docs.openstack.org/image-guide/image-formats.html) for more information.

* The public network should be created by the OSP admin. Verify the name/ID of the 'External' network:
```
openstack network list --long -c ID -c Name -c "Router Type"
+--------------------------------------+----------------+-------------+
| ID                                   | Name           | Router Type |
+--------------------------------------+----------------+-------------+
| 148a8023-62a7-4672-b018-003462f8d7dc | public_network | External    |
+--------------------------------------+----------------+-------------+
```

## Current Expected Behavior

As mentioned, OpenStack support is still experimental. Currently:

* Deploys an isolated tenant network
* Deploys a instance used as a 'service VM' that hosts a load balancer for the
OpenShift API and as an internal DNS for the instances
* Deploys a bootstrap instance to bootstrap the OpenShift cluster
* Deploys 3 master nodes
* Once the masters are deployed, the bootstrap instance is destroyed

The installer fails to end gracefully as the openshift-console is not deployed
because there are no nodes available.

**NOTE** The worker nodes are still a WIP

### Workarounds

#### External DNS

While deploying the cluster, the installer will hang trying to reach the API as
the node running the installer cannot resolve the service VM (the cluster
should still come up successfully within the isolated network). As a temporary
workaround you can add the service VM floating IP and hostname to the
`/etc/hosts` file as:

```
$ cat /etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
10.19.115.117 <cluster-name>-api.<domain>
```

If you do expose the cluster, the installer should make it far enough along to
bring up the HA control plane and tear down the bootstrap node.  It will then
hang waiting for the console to come up.

```
DEBUG Still waiting for the console route: the server is currently unable to
handle the request (get routes.route.openshift.io)
...
FATAL waiting for openshift-console URL: context deadline exceeded
```

#### Create the openstack-credentials Secret

This is necessary for creating the worker nodes and scaling the cluster. The
actuator needs to have access to the OpenStack credentials. They are being read
from an Kubernetes Secret object.

This is a post-deployment operation: it should be done after the bootstrap node
has been removed by the installer.

1. Create a file called `secret.yaml` with the following contents:

```
$ cat << EOF > secret.yaml
apiVersion: v1
data:
  clouds.yaml: $(cat $HOME/.config/openstack/clouds.yaml | base64 -w0)
kind: Secret
metadata:
  name: openstack-credentials
type: Opaque
EOF
```

2. Add the Secret to the cluster:

```
$ oc create -n openshift-cluster-api -f secret.yaml
```

## Using an External Load Balancer

This documents how to shift from the API VM load balancer, which is
intended for initial cluster deployment and not highly available, to an
external load balancer.

The load balancer must serve ports 8443, 443, and 80 to any users of
the system.  Port 49500 is for serving ignition startup configurations
to the OpenShift nodes and should not be reachable outside of the cluster.

The first step is to add floating IPs to all the master nodes:

* `openstack floating ip create --port master-port-0 <public network>`
* `openstack floating ip create --port master-port-1 <public network>`
* `openstack floating ip create --port master-port-2 <public network>`

Once complete you can see your floating IPs using:

* `openstack server list`

These floating IPs can then be used by the load balancer to access
the cluster.  An example haproxy configuration for port 8443 is below.
The other port configurations are identical.

```
listen <cluster name>-api-8443
    bind 0.0.0.0:8443
    mode tcp
    balance roundrobin
    server <cluster name>-master-2 <floating ip>:8443 check
    server <cluster name>-master-0 <floating ip>:8443 check
    server <cluster name>-master-1 <floating ip>:8443 check
```

The next step is to allow network access from the load balancer network
to the master nodes:

* `openstack security group rule create master --remote-ip <load balancer CIDR> --ingress --protocol tcp --dst-port 8443`
* `openstack security group rule create master --remote-ip <load balancer CIDR> --ingress --protocol tcp --dst-port 443`
* `openstack security group rule create master --remote-ip <load balancer CIDR> --ingress --protocol tcp --dst-port 80`

You could also specify a specific IP address with /32 if you wish.

You can verify the operation of the load balancer now if you wish, using the
curl commands given below.

Now the DNS entry for <cluster name>-api.<base domain> needs to be updated
to point to the new load balancer:

* `<load balancer ip> <cluster-name>-api.<base domain>`

The external load balancer should now be operational along with your own
DNS solution. It's best to test this configuration before removing
the API. The following curl command is an example of how
to check functionality:

`curl https://<loadbalancer-ip>:8443/version --insecure`

Result:

```json
{
  "major": "1",
  "minor": "11+",
  "gitVersion": "v1.11.0+ad103ed",
  "gitCommit": "ad103ed",
  "gitTreeState": "clean",
  "buildDate": "2019-01-09T06:44:10Z",
  "goVersion": "go1.10.3",
  "compiler": "gc",
  "platform": "linux/amd64"
}
```

Another useful thing to check is that the ignition configurations are
only available from within the deployment. The following command should
only succeed from a node in the OpenShift cluster:

* `curl https://<loadbalancer ip>:49500/config/master --insecure`

Now that the DNS and load balancer has been moved, we can take down the existing
api VM:

* `openstack server delete <cluster name>-api`

## Troubleshooting

See the [troubleshooting installer issues in OpenStack](./troubleshooting.md) guide.

## Reporting Issues

Please see the [Issue Tracker][issues_openstack] for current known issues.
Please report a new issue if you do not find an issue related to any trouble
you’re having.

[issues_openstack]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+openstack
