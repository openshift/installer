# Limits

You can find a comprehensive list of the default AWS service limits published here:

[AWS Service Limits][service-limits]

Below, we'll identify OpenShift cluster needs and how those impact some of those limits.
We will also use [aws cli][aws-cli] to try and explan
where limits are set (viewable) and how to show your current usage of resources.

Many commands below are based on [Trusted Advisor][trusted-advisor] and often this can be used to print both the limit and 
the usage for a service. However to check those limits/usage you need to first identify a check_id. 
These check_id's can be pulled from:

```console
$ aws support describe-trusted-advisor-checks --language en --query "checks[].{id: @.id, name: @.name, category: @.category}" --output text
```

Should you need to increase a limit (with any of the following) you need to open a ticket with AWS support, according to
[AWS Service Limits][service-limits].

### Example: Using N. Virginia (us-east-1) to increase a limit for a VPC Elastic IP

To use N. Virginia (us-east-1) for a new cluster, please submit a limit increase for VPC Elastic IPs similar to the
following in the support dashboard (to create more than one cluster, a higher limit will be necessary):

![Increase Elastic IP limit in AWS](images/support_increase_elastic_ip.png)

## S3

There is a default limit of 100 S3 buckets per account. The installation creates a bucket temporarily. Also, the
registry component creates a permanent bucket. This will limit the number of clusters per account to 99 initially. To
support additional clusters, you must open a support case with AWS.

You can not view the limits for S3 buckets, however you can view how many buckets you are using (what the name is), 
by running the following:

```console
$ aws s3api --query 'join(`\n`, @.Buckets[].Name)' --output text list-buckets
```

## VPC

Each cluster creates its own VPC. The default limit of VPCs per region is 5 and will allow 5 clusters. To have more
than 5 clusters, you will need to increase this limit.

To see both the limit and the usage for VPC's you can use:

```console
$ aws support describe-trusted-advisor-check-result --check-id jL7PP0l7J9 --query 'result.flaggedResources[].{service_limit: @.metadata[2], region:@.region, status: @.status, limit: @.metadata[3], usage: @.metadata[4]}' --output text
```

Each cluster also creates a VPC Gateway Endpoint for a stable connection to S3. The default limit of VPC Gateway
Endpoints per region is 20 and will allow 20 clusters. To have more than 20 clusters, you will need to increase this
limit.

You can not view the limits for VPC Endpoints, however you can view how many endpoings you are using (what the name and
creation date and type are), by running the following:

- Be sure to change the region (from us-east-1) if you are using a different region.

```console
$ aws --region us-east-1 ec2 describe-vpc-endpoints --query "VpcEndpoints[].{id:@.VpcEndpointId, type: @.VpcEndpointType, connected_VPC: @.VpcId, connected_service: @.ServiceName, creation_timestamp: @.CreationTimestamp, state: @.State}" --output text
```

## Elastic Network Interfaces (ENI)

The default installation creates 21 + the number of availability zones of ENIs (e.g. 21 + 3 = 24 ENIs for a three-zone cluster).
The default limit per region is 350. Additional ENIs are created for additional machines and elastic load balancers
created by cluster usage and deployed workloads. A service limit increase here may be required to satisfy the needs of
additional clusters and deployed workloads.

You can not view the limits for Elastic Network Interfaces, however you can view how many interfaces you are using
(what the name and type are), by running the following:

- Be sure to change the region (from us-east-1) if you are using a different region.

```console
$ aws ec2 --region us-east-1 describe-network-interfaces --query "NetworkInterfaces[].{associated_vpc: @.VpcId, associated_subnet: @.SubnetId, zone: @.AvailabilityZone, status: @.Status, id: @.NetworkInterfaceId}" --output text
```

## Elastic IP (EIP)

By default, the installer distributes control-plane and compute machines across [all availability zones within a region][availability-zones] to provision the cluster in a highly available configuration.
Please see [this map][az-map] for a current region map with availability zone count.
We recommend selecting regions with 3 or more availability zones.
You can [provide an install-config](../overview.md#multiple-invocations) to [configure](customization.md) the installer to use specific zones to override that default.

The installer creates a public and private subnet for each configured availability zone.
In each private subnet, a separate [NAT Gateway][nat-gateways] is created and requires a separate [EC2-VPC Elastic IP (EIP)][elastic-ip].
The default limit of 5 is sufficient for a single cluster, unless you have configured your cluster to use more than five zones.
For multiple clusters, a higher limit will likely be required (and will certainly be required to support more than five clusters, even if they are each single-zone clusters).

To see both the limit and the usage for EIP's you can use:

```console
$ aws support describe-trusted-advisor-check-result --check-id lN7RR0l7J9 --query 'result.flaggedResources[].{service_limit: @.metadata[2], region:@.region, status: @.status, limit: @.metadata[3], usage: @.metadata[4]}' --output text
```

## NAT Gateway

The default limit for NAT Gateways is 5 per availability zone. This is sufficient for up to 5 clusters in a dedicated
account. If you intend to create more than 5 clusters, you will need to request an increase to this limit.

You can not view the limits for VPC Endpoints, however you can view how many endpoings you are using (what the name and
creation date and type are), by running the following:

- Be sure to change the region (from us-east-1) if you are using a different region.

```console
$ aws --region us-east-1 ec2 describe-nat-gateways --query "NatGateways[].{id: @.NatGatewayId, state: @.State, associated_subnet_id: @.SubnetId, associated_vpc_id: @.VpcId}" --output text
```

## VPC Gateway

The default limit of VPC Gateways (for S3 access) is 20. Each cluster will create a single S3 gateway endpoint within
the new VPC. If you intend to create more than 20 clusters, you will need to request an increase to this limit.

To see both the limit and the usage for EIP's you can use:

```console
$ aws support describe-trusted-advisor-check-result --check-id kM7QQ0l7J9 --query 'result.flaggedResources[].{service_limit: @.metadata[2], region:@.region, status: @.status, limit: @.metadata[3], usage: @.metadata[4]}' --output text
```

## Security Groups

Each cluster creates distinct security groups. The default limit of 2,500 for new accounts allows for many clusters
to be created. The security groups which exist after the default install are:

  1. VPC default
  1. Master
  1. Worker
  1. Router/Ingress

You can not view the limits for Security Groups, however you can view how many Security Groups you have defined using,
by running the following:

```console
$ aws ec2 describe-security-groups --query "SecurityGroups[].{id: @.GroupId, name: @.GroupName}" --output text
```

## Instance Limits

By default, a cluster will create:

* One m4.large bootstrap machine (removed after install)
* Three m4.xlarge master nodes.
* Three m4.large worker nodes.

Currently, these instance type counts are within a new account's default limit.
If you intend to start with a higher number of workers, enable autoscaling and large workloads
or a different instance type, please ensure you have the necessary remaining instance count within the instance type's
limit to satisfy the need. If not, please ask AWS to increase the limit via a support case.

You can then confirm the limits and usage (for on Demand Instances) by running the following:

```console
$ aws support describe-trusted-advisor-check-result --check-id 0Xc6LMYG8P --query 'result.flaggedResources[].{service_limit: @.metadata[2], region:@.region, status: @.status, limit: @.metadata[3], usage: @.metadata[4]}' --output text
```

## Elastic Load Balancing (ELB/NLB)

By default, each cluster will create 2 network load balancers for the master API server (1 internal, 1 external) and a
single classic elastic load balancer for the router. Additional Kubernetes LoadBalancer Service objects will create
additional [load balancers][load-balancing].

You can then confirm the limits and usage (for Classic Load Balancers) by running the following:

```console
$ aws support describe-trusted-advisor-check-result --check-id iK7OO0l7J9 --query 'result.flaggedResources[].{service_limit: @.metadata[2], region:@.region, status: @.status, limit: @.metadata[3], usage: @.metadata[4]}' --output text
```

And for	Network	Load Balancers:

- Be sure to change the region (from us-east-1) if you are using a different region.

```sh
$ aws --region us-east-1 elbv2 describe-account-limits --query "Limits[? @.Name == 'network-load-balancers'].Max" --output text
$ aws --region us-east-1 elbv2 describe-load-balancers --query 'join(`\n`, @.LoadBalancers[].LoadBalancerArn)' --output text | wc -l
```

[availability-zones]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
[az-map]: https://aws.amazon.com/about-aws/global-infrastructure/
[elastic-ip]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/elastic-ip-addresses-eip.html
[load-balancing]: https://aws.amazon.com/elasticloadbalancing/
[nat-gateways]: https://docs.aws.amazon.com/vpc/latest/userguide/vpc-nat-gateway.html
[service-limits]: https://docs.aws.amazon.com/general/latest/gr/aws_service_limits.html
[aws-cli]: https://aws.amazon.com/cli/
[trusted-advisor]: https://aws.amazon.com/premiumsupport/technology/trusted-advisor/
