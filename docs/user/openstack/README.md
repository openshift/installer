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

* Swift must be enabled.  The user must have `swiftoperator` permissions and
  `temp-url` support must be enabled.
  * `openstack role add --user <user> --project <project> swiftoperator`
  * `openstack object store account set --property Temp-URL-Key=superkey`

* You may need to increase the security group related quotas from their default
  values.
  * For example: `openstack quota set --secgroups 100 --secgroup-rules 1000
    openshift`

## Current Expected Behavior

As mentioned, OpenStack support is still experimental.  The installer will
launch a cluster on an isolated tenant network.  To access your cluster, you
can create a floating IP address and assign it to the load balancer service VM.

* `openstack floating ip create ${EXTERNAL_NETWORK}`
* `openstack floating ip set --port lb-port ${FLOATING_IP_ADDRESS}`

The service VM also hosts a DNS server that has enough records to bring up the
cluster.  If you add the `${FLOATING_IP_ADDRESS}` as your first `nameserver`
entry in `/etc/resolv.conf`, the installer will be able to look up the address
needed to reach the API.

If you don’t expose the cluster and add a hosts entry, the installer will hang
trying to reach the API.  However, the cluster should still come up
successfully within the isolated network.

If you do expose the cluster, the installer should make it far enough along to
bring up the HA control plane and tear down the bootstrap node.  It will then
hang waiting for the console to come up.

`DEBUG Still waiting for the console route: the server is currently unable to
handle the request (get routes.route.openshift.io)`

## Reporting Issues

Please see the [Issue Tracker][issues_openstack] for current known issues.
Please report a new issue if you do not find an issue related to any trouble
you’re having.

[issues_openstack]: https://github.com/openshift/installer/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+openstack
