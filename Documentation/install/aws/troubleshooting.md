# Troubleshooting Tectonic Installer on AWS

## Tectonic Installer stalls if the AWS security token expires

Tectonic installation on AWS fails if the AWS Security Token Service (STS) token expires before installation is complete. AWS imposes a hard 1 hour limit on the token. If your token expires during an installation, refresh the credentials, and destroy the partially created cluster on AWS. Before starting again, refresh the token, then reinstall the cluster using the same tfvars or installer settings as the original install.

## Route 53 DNS resolution

### DNS caching

An issue arises when a domain's Address record (A record) resolution is attempted before Route 53 publishes the cluster's A record and the NXDOMAIN response is cached in the NCACHE (RFC2308). This negative response may be cached for up to the number of seconds set in the domain's SOA record's TTL. Resolution fails until the negative caching TTL expires. These TTLs are typically large enough to disrupt the installation. The current workaround is to ensure your TTLs are set to a low interval, or to wait for them to expire, then proceed with the installation.

### Internal VPC, ELB and/or Hosted Zone

This can be applicable to a vareity of situations which restrict external access to your Tectonic cluster.

* Internal ELB: `tectonic_aws_external_vpc_public = false`
* Public Route53 zone is not delegated from a public hosted zone
* Custom Security Groups or Subnet ACLs which restrict access from outside the VPC to ELB
* Custom Route Tables which do not provide egress via Internet Gateway or NAT Gateway

The above list represents the most common configuration choices which have the potential to interfere with external access to the VPC.

The simplest way to install Tectonic in this situation is to run the installer on a "bastion" CoreOS EC2 instance within the same existing VPC or existing subnet that you are installing to. After the cluster install is complete and you have your assets.zip file safely stored, you can terminate the bastion host. It's recommended that you use a `t2.small` or larger instance for the bastion host.

Whether you are running the Tectonic Installer on a bastion host within the VPC or externally via a VPN connection, the following troubleshooting items can be relevant.

 * __The installer cannot resolve DNS records for the tectonic ingress load balancer__:

    * _If you want to use public DNS to resolve Tectonic ingress_: In order for a Route53 public hosted zone to be resolvable outside the VPC, it must be delegated to. [More information on delegating to Route53 subdomains](http://docs.aws.amazon.com/Route53/latest/DeveloperGuide/CreatingNewSubdomain.html#UpdateDNSParentDomain). Make sure the public hosted zone used by the Tectonic has proper delegation set up before attempting an install.

    * _If you want to confine Route53 records to only be resolvable within the VPC_: Delegation from a public DNS zone is not necessary.

    Suggested solutions:

      * Run the tectonic installer on an EC2 instance within that VPC, and make sure the `dhcp-options-set` for that machine specifies the VPC's DNS server

      * If you are using a VPN connection to to allow the installer access to the VPC, make sure the machine hosting the Tectonic installer is using the VPC's DNS endpoint as the primary DNS server.

    From the machine hosting the Tectonic installer, `dig <tectonic_cluster_name>.<tectonic_base_domain> NS +short` should yield the same list of NS servers you see listed for the applicable public hosted zone via the AWS console or API.

 * __"no route to host" when attempting to establish an HTTP connection with tectonic ingress load balancer__

   Things to check:

     * Is `tectonic_aws_external_vpc_public = false`? This means that the ingress ELB is only accessible from within the VPC.

       In this case, you'll either need to run the installer on an EC2 instance within the VPC or establish a VPN connection into the VPC.

     * Route tables for either EC2 instance running tectonic installer or Virtual Private Gateway providing VPN connection from installer to AWS VPC.

 * __"connection refused" when attempting to establish an HTTP connection with tectonic ingress load balancer__

   This can indicate a security group rule and/or subnet ACL which is preventing the installer from establishing TCP connection with the ELB.

When debugging an install that failed for the above reasons, you can use `traceroute` to determine if the installer has proper access to the Tectonic ingress ELB.

On the machine hosting the Tectonic installer, the following command should succeed (assuming the ELB and Route53 DNS records have already been created):

```sh
traceroute -T -p 443 <tectonic_cluster_name>.<tectonic_base_domain>
```

This validates that the installer is able to establish a TCP connection with the Tectonic ingress ELB. This will validate DNS, IP routing and firewall rules between the installer and the Tectonic ingress ELB.

## Domain name can't be changed

The domain configured for Route 53 name service and the domain names selected for Tectonic and Controller DNS names during install cannot be easily changed later. If a cluster's domain name must change, set up a new cluster with the new domain name and migrate cluster work to it.

## VPC peering

Pod, Service, and Instance CIDRs may overlap with ranges of [peered VPCs](http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-peering.html), but only if the `pcx` interface is *not* in the route table attached to the Kubernetes subnets. Be certain to configure these CIDRs appropriately.

## Resource tagging

Most resources in AWS support tagging, and the Tectonic installer tags as many of these resources as possible. Unfortunately certain resources are created implicitly by others, and no API exists to tag them, ex. EC2 instances create network interfaces that cannot be tagged by the `aws_instance` resource. The following notes describe which resources are not or are incompletely tagged.

* Autoscaling Group-controlled instances are not tagged with user-defined tags, only with defaults
* Autoscaling Group-controlled instance volumes are not tagged with any tags
* Default EC2 Network ACLs are not tagged with any tags
* EC2 Network Interfaces implicitly created by EC2 Instances are not tagged with any tags

## Community Support Forum

Make sure to check out the [community support forum](https://github.com/coreos/tectonic-forum/issues) to work through issues, report bugs, identify documentation requirements, or add feature requests.
