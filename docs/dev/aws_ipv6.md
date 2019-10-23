# AWS IPv6 Dev and Test Environment

The installer includes code to enable an experimental IPv6 dev and test
environment on AWS.  This is off by default and currently only intended for
those working on enabling IPv6 in OpenShift.

## Enabling the Environment

To enable IPv6 in your AWS environment, set the following environment variable
before running the installer:

```bash
    export TF_VAR_aws_use_ipv6=”true”
```

## AWS Network Environment

AWS does not support a single-stack (IPv6 only) environment, but it does
support a dual-stack (IPv4 and IPv6) environment, so that’s what is enabled
here.  This is a summary of the changes to the network environment:

* The VPC has IPv6 enabled and a `/56` IPv6 CIDR will be allocated by AWS.
* Each Subnet will have an IPv6 `/64` subnet allocated to it.
* All IPv4 specific security group rules have corresponding IPv6 rules created.
* AWS Network Load Balancers (NLBs) do not support IPv6, so classic Elasic Load
  Balancers (ELBs) are created, instead.
* The IPv4 NLBs are disabled.
* IPv6 DNS records (AAAA) are created and the IPv4 (A) records are disabled.
* IPv6 routing is configured.  Since all instances get global IPv6 addresses,
  NAT is not used from the instances out to the internet.

## Install Configuration

Here is the suggested network configuration for `install-config.yaml`:

```yaml
networking:
  clusterNetwork:
  - cidr: fd01::/48
    hostPrefix: 64
  machineCIDR: 10.0.0.0/16
  networkType: OVNKubernetes
  serviceNetwork:
  - fd02::/112
```

Note that an IPv4 CIDR is still used for `machineCIDR` since AWS will provide a
dual-stack (IPv4 and IPv6) environment.  We must specify the IPv4 CIDR and AWS
will automatically allocate an IPv6 CIDR.

## Current Status of IPv6

Note that IPv6 support is under heavy development across many components in
OpenShift, so the use of some custom images may be needed to include fixes to
known issues.  Coordination of work-in-progress is out of scope for this
document.
