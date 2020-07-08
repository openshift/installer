# Limits

You can find a comprehensive list of the default AWS service limits published here:

[AWS Service Limits][service-limits]

Below, we'll identify OpenShift cluster needs and how those impact some of those limits.

## S3

There is a default limit of 100 S3 buckets per account. The installation creates a bucket temporarily. Also, the
registry component creates a permanent bucket. This will limit the number of clusters per account to 99 initially. To
support additional clusters, you must open a support case with AWS.

## VPC

Each cluster creates its own VPC. The default limit of VPCs per region is 5 and will allow 5 clusters. To have more
than 5 clusters, you will need to increase this limit.

Each cluster also creates a VPC Gateway Endpoint for a stable connection to S3. The default limit of VPC Gateway 
Endpoints per region is 20 and will allow 20 clusters. To have more than 20 clusters, you will need to increase this 
limit.

## Elastic Network Interfaces (ENI)

The default installation creates 21 + the number of availability zones of ENIs (e.g. 21 + 3 = 24 ENIs for a three-zone cluster).
The default limit per region is 350. Additional ENIs are created for additional machines and elastic load balancers
created by cluster usage and deployed workloads. A service limit increase here may be required to satisfy the needs of
additional clusters and deployed workloads.

## Elastic IP (EIP)

By default, the installer distributes control-plane and compute machines across [all availability zones within a region][availability-zones] to provision the cluster in a highly available configuration.
Please see [this map][az-map] for a current region map with availability zone count.
We recommend selecting regions with 3 or more availability zones.
You can [provide an install-config](../overview.md#multiple-invocations) to [configure](customization.md) the installer to use specific zones to override that default.

The installer creates a public and private subnet for each configured availability zone.
In each private subnet, a separate [NAT Gateway][nat-gateways] is created and requires a separate [EC2-VPC Elastic IP (EIP)][elastic-ip].
The default limit of 5 is sufficient for a single cluster, unless you have configured your cluster to use more than five zones.
For multiple clusters, a higher limit will likely be required (and will certainly be required to support more than five clusters, even if they are each single-zone clusters).

### Example: Using North Virginia (us-east-1)

North Virginia (us-east-1) has six availablity zones, so a higher limit is required unless you configure your cluster to use fewer zones.
To support the default, all-zone installation, please submit a limit increase for VPC Elastic IPs similar to the following in the support dashboard (to create more than one cluster, a higher limit will be necessary):

![Increase Elastic IP limit in AWS](images/support_increase_elastic_ip.png)

## NAT Gateway

The default limit for NAT Gateways is 5 per availability zone. This is sufficient for up to 5 clusters in a dedicated
account. If you intend to create more than 5 clusters, you will need to request an increase to this limit.

## VPC Gateway

The default limit of VPC Gateways (for S3 access) is 20. Each cluster will create a single S3 gateway endpoint within
the new VPC. If you intend to create more than 20 clusters, you will need to request an increase to this limit.

## Security Groups

Each cluster creates distinct security groups. The default limit of 2,500 for new accounts allows for many clusters
to be created. The security groups which exist after the default install are:

  1. VPC default
  1. Master
  1. Worker
  1. Router/Ingress

## vCPU Limits

By default, a cluster will create:

* One m4.large bootstrap machine (2 vCPUs but removed after install)
* Three m5.xlarge master nodes (4 vCPUs each).
* Three m5.large worker nodes (2 vCPUs each).

Currently, these vCPU counts are not within a new account's default limit. The default limit is 1 but for all these instances you will need 20. To increase the limit you have to [contact the AWS support](https://console.aws.amazon.com/support/cases?#/create?issueType=service-limit-increase&limitType=ec2-instances).
If you intend to start with a higher number of workers, enable autoscaling and large workloads
or a different instance type, please ensure you have the necessary remaining vCPU count within the vCPU
limit to satisfy the need. To calculate the vCPU limit you can use the limits calculator in the EC2 console (EC2 -> Limits -> Limits calculator).

## Elastic Load Balancing (ELB/NLB)

By default, each cluster will create 2 network load balancers for the master API server (1 internal, 1 external) and a
single classic elastic load balancer for the router. Additional Kubernetes LoadBalancer Service objects will create
additional [load balancers][load-balancing].

[availability-zones]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
[az-map]: https://aws.amazon.com/about-aws/global-infrastructure/
[elastic-ip]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/elastic-ip-addresses-eip.html
[load-balancing]: https://aws.amazon.com/elasticloadbalancing/
[nat-gateways]: https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html
[service-limits]: https://docs.aws.amazon.com/general/latest/gr/aws_service_limits.html
