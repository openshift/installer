# Troubleshooting Tectonic Installer on AWS

## Route 53 DNS resolution

An issue arises when a domain's Address record (A record) resolution is attempted before Route 53 publishes the cluster's A record and the NXDOMAIN response is cached in the NCACHE (RFC2308). This negative response may be cached for up to the number of seconds set in the domain's SOA record's TTL. Resolution fails until the negative caching TTL expires. These TTLs are typically large enough to disrupt the installation. The current workaround is to ensure your TTLs are set to a low interval, or to wait for them to expire, then proceed with the installation.

## Domain name can't be changed

The domain configured for Route 53 name service and the domain names selected for Tectonic and Controller DNS names during install cannot be easily changed later. If a cluster's domain name must change, set up a new cluster with the new domain name and migrate cluster work to it.

## VPC peering

Pod, Service, and Instance CIDRs may overlap with ranges of [peered VPCs](http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-peering.html), but only if the `pcx` interface is *not* in the route table attached to the Kubernetes subnets. Be certain to configure these CIDRs appropriately.

## Community Support Forum

Make sure to check out the [community support forum](https://github.com/coreos/tectonic-forum/issues) to work through issues, report bugs, identify documentation requirements, or add feature requests.
