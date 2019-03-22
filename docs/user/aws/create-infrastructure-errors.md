# Errors with Infrastructure Creation

This document describes errors with infrastructure creation on AWS.
For more troubleshooting suggestions, look [here](../troubleshooting.md).

## Instances time out with `pending`

The error looks like:

```
aws_instance.master.1: Error launching source instance: timeout while waiting for state to become 'success' (timeout: 30s)
```

or:

```
aws_instance.master.1: Error waiting for instance (i-0eb201c7376b6f94e) to become ready: timeout while waiting for state to become 'running' (last state: 'pending', timeout: 10m0s)"
```

This happens when AWS does not have the capacity in the target availability zone to fulfill the instance request.
There may be multiple datacenters within a given availability zone, and the limitation may be restricted to a subset of datacenters within the affected zone.
You probably want to destroy the failed cluster, after which you have a few choices:

* You can try a fresh install with the same parameters, hoping that the capacity has become available or that you happen to be assigned to an unaffected datacenter.
* You can configure a different region or availability zone(s) to find sufficient capacity.
* You can configure a different instance type to find sufficient capacity.
