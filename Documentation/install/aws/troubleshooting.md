# Troubleshooting Tectonic Installer on AWS

## Tectonic Installer stalls if the STS security token expires during a session

Tectonic installation on AWS fails if the AWS Security Token Service (STS) token expires before installation is complete. This issue arises when multi-factor authentication is configured for the IAM role. If the installation is resumed with a progress file, resource names might conflict because the cluster has been partially created on AWS. To work around, delete the AWS CloudFormation stack before using the progress file to restart the install. In the AWS console, click the *CloudFormation* option under *Management Tools*. Delete the stack with the name of the cluster whose installation stalled.

## Route53 DNS resolution

An issue arises when a domain's Address record (A record) resolution is attempted before Route53 publishes the cluster's A record and the NXDOMAIN response is cached in the NCACHE (RFC2308). This negative response may be cached for up to the number of seconds set in the domain's SOA record's TTL. Resolution fails until the negative caching TTL expires. These TTLs are typically large enough to disrupt the installation. The current workaround is to ensure your TTLs are set to a low interval, or to wait for them to expire, then proceed with the installation.

## Domain name can't be changed

The domain configured for Route53 name service and the domain names selected for Tectonic and Controller DNS names during install cannot be easily changed later. If a cluster's domain name needs to change, set up a new cluster with the new domain name and migrate cluster work to it.

## VPC peering not supported

A VPC containing a Tectonic cluster should not be [peered](http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-peering.html) with other VPCs due to the [potential for IP routing conflicts](http://ben.straub.cc/2015/08/19/kubernetes-aws-vpc-peering/).

## Community Support Forum

Make sure to check out the [community support forum](https://github.com/coreos/tectonic-forum/issues) to work through issues, report bugs, identify documentation requirements, or put in feature requests.
