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
as `raw` or `qcow2`. See [https://docs.openstack.org/image-guide/image-formats.html](Disk and container formats for images) for more information.

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

## Reporting Issues

Please see the [Issue Tracker][issues_openstack] for current known issues.
Please report a new issue if you do not find an issue related to any trouble
you’re having.

[issues_openstack]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+openstack
