# Install: User Provided Infrastructure (UPI)

The steps for performing a UPI-based install are outlined here. Several [CloudFormation][cloudformation] templates are
provided to assist in completing these steps or to help model your own.  You are also free to create the required
resources through other methods; the CloudFormation templates are just an example.

## Create Configuration

Create an install configuration as for [the usual approach](install.md#create-configuration):

```console
$ openshift-install create install-config
? SSH Public Key /home/user_id/.ssh/id_rsa.pub
? Platform aws
? Region us-east-2
? Base Domain example.com
? Cluster Name openshift
? Pull Secret [? for help]
```

## Edit Manifests

Use [a staged install](../overview.md#multiple-invocations) to make some adjustments which are not exposed via the install configuration.

```console
$ openshift-install create manifests
INFO Consuming "Install Config" from target directory
```

### Remove Machines and MachineSets

Remove the control-plane Machines and compute MachineSets, because we'll be providing those ourselves and don't want to involve [the machine-API operator][machine-api-operator]:

```console
$ rm -f openshift/99_openshift-cluster-api_master-machines-*.yaml openshift/99_openshift-cluster-api_worker-machinesets-*.yaml
```

You are free to leave the compute MachineSets in if you want to create compute machines via the machine API, but if you do you may need to update the various references (`subnet`, etc.) to match your environment.

### Remove DNS Zones

If you don't want [the ingress operator][ingress-operator] to create DNS records on your behalf, remove the `privateZone` and `publicZone` sections from the DNS configuration:

```sh
python -c '
import yaml;
path = "manifests/cluster-dns-02-config.yml";
data = yaml.load(open(path));
del data["spec"]["publicZone"];
del data["spec"]["privateZone"];
open(path, "w").write(yaml.dump(data, default_flow_style=False))'
```

If you do so, you'll need to [add ingress DNS records manually](#add-the-ingress-dns-records) later on.

## Create Ignition Configs

Now we can create the bootstrap Ignition configs:

```console
$ openshift-install create ignition-configs
```

After running the command, several files will be available in the directory.

```console
$ tree
.
├── auth
│   └── kubeconfig
├── bootstrap.ign
├── master.ign
├── metadata.json
└── worker.ign
```

### Extract Infrastructure Name from Ignition Metadata

Many of the operators and functions within OpenShift rely on tagging AWS resources. By default, Ignition
generates a unique cluster identifier comprised of the cluster name specified during the invocation of the installer
and a short string known internally as the infrastructure name. These values are seeded in the initial manifests within
the Ignition configuration. To use the output of the default, generated 
`ignition-configs` extracting the internal infrastructure name is necessary.

An example of a way to get this is below: 

```
$ jq -r .infraID metadata.json 
openshift-vw9j6
```

## Create/Identify the VPC to be Used

You may create a VPC with various desirable characteristics for your situation (VPN, route tables, etc.). The
VPC configuration and a CloudFormation template is provided [here](../../../upi/aws/cloudformation/01_vpc.yaml).

A created VPC via the template or manually should approximate a setup similar to this:

<div style="text-align:center">
  <img src="images/install_upi_vpc.svg" width="100%" />
</div>

## Create DNS entries and Load Balancers for Control Plane Components

The DNS and load balancer configuration within a CloudFormation template is provided
[here](../../../upi/aws/cloudformation/02_cluster_infra.yaml). It uses a public hosted zone and creates a private hosted
zone similar to the IPI installation method. 
It also creates load balancers, listeners, as well as hosted zone and subnet tags the same way as the IPI
installation method. 
This template can be run multiple times within a single VPC and in combination with the VPC
template provided above.

### Optional: Manually Create Load Balancer Configuration

It is needed to create a TCP load balancer for ports 6443 (the Kubernetes API and its extensions) and 22623 (Ignition
configurations for new machines).  The targets will be the master nodes.  Port 6443 must be accessible to both clients
external to the cluster and nodes within the cluster. Port 22623 must be accessible to nodes within the cluster.

### Optional: Manually Create Route53 Hosted Zones & Records

For the cluster name identified earlier in [Create Ignition Configs](#create-ignition-configs), you must create a DNS entry which resolves to your created load balancer.
The entry `api.$clustername.$domain` should point to the external load balancer and `api-int.$clustername.$domain` should point to the internal load balancer.

## Create Security Groups and IAM Roles

The security group and IAM configuration within a CloudFormation template is provided
[here](../../../upi/aws/cloudformation/03_cluster_security.yaml). Run this template to get the minimal and permanent
set of security groups and IAM roles needed for an operational cluster. It can also be inspected for the current
set of required rules to facilitate manual creation.

## Launch Temporary Bootstrap Resource

The bootstrap launch and other necessary, temporary security group plus IAM configuration and a CloudFormation
template is provided [here](../../../upi/aws/cloudformation/04_cluster_bootstrap.yaml). Upload your generated `bootstrap.ign`
file to an S3 bucket in your account and run this template to get a bootstrap node along with a predictable clean up of
the resources when complete. It can also be inspected for the set of required attributes via manual creation.

## Launch Permanent Master Nodes

The master launch and other necessary DNS entries for etcd are provided within a CloudFormation
template [here](../../../upi/aws/cloudformation/05_cluster_master_nodes.yaml). Run this template to get three master
nodes. It can also be inspected for the set of required attributes needed for manual creation of the nodes, DNS entries
and load balancer configuration.

## Monitor for `bootstrap-complete` and Initialization

```console
$ bin/openshift-install wait-for bootstrap-complete
INFO Waiting up to 30m0s for the Kubernetes API at https://api.test.example.com:6443...
INFO API v1.12.4+c53f462 up
INFO Waiting up to 30m0s for the bootstrap-complete event...
```

## Destroy Bootstrap Resources

At this point, you should delete the bootstrap resources. If using the CloudFormation template, you would [delete the
stack][delete-stack] created for the bootstrap to clean up all the temporary resources.

## Launch Additional Compute Nodes

You may create compute nodes by launching individual EC2 instances discretely or by automated processes outside the cluster (e.g. Auto Scaling Groups).
You can also take advantage of the built in cluster scaling mechanisms and the machine API in OpenShift, as mentioned [above](#create-ignition-configs).
In this example, we'll manually launch instances via the CloudFormatio template [here](../../../upi/aws/cloudformation/06_cluster_worker_node.yaml).
You can launch a CloudFormation stack to manage each individual compute node (you should launch at least two for a high-availability ingress router).
A similar launch configuration could be used by outside automation or AWS auto scaling groups.

#### Approving the CSR requests for nodes

The CSR requests for client and server certificates for nodes joining the cluster will need to be approved by the administrator.
You can view them with:

```console
$ oc get csr
NAME        AGE     REQUESTOR                                                                   CONDITION
csr-8b2br   15m     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-8vnps   15m     system:serviceaccount:openshift-machine-config-operator:node-bootstrapper   Approved,Issued
csr-b96j4   25s     system:node:ip-10-0-52-215.us-east-2.compute.internal                       Approved,Issued
csr-bfd72   5m26s   system:node:ip-10-0-50-126.us-east-2.compute.internal                       Pending
csr-c57lv   5m26s   system:node:ip-10-0-95-157.us-east-2.compute.internal                       Pending
...
```

Administrators should carefully examine each CSR request and approve only the ones that belong to the nodes created by them.
CSRs can be approved by name, for example:

```sh
oc adm certificate approve csr-bfd72
```

## Configure Router for UPI

The Ingress operator manages DNS and LoadBalancers. It makes use of tags on HostedZones to identify which public and 
private zones are to be updated from the cluster by the operator as objects are created in the cluster. It makes use
of tags on subnets to identify those to associate with Service objects of type LoadBalancer created in the cluster.

The tags used for finding HostedZones used by the operator
are fulfilled by the CloudFormation template [here](../../../upi/aws/cloudformation/02_cluster_infra.yaml).

An example of the `spec` for DNS configuration is below:

```
$ oc get dns -o yaml
apiVersion: v1
items:
- apiVersion: config.openshift.io/v1
  kind: DNS
  metadata:
    creationTimestamp: 2019-03-28T12:31:10Z
    generation: 1
    name: cluster
    namespace: ""
    resourceVersion: "395"
    selfLink: /apis/config.openshift.io/v1/dnses/cluster
    uid: 5e51dd25-5155-11e9-befc-02d75ce1a902
  spec:
    baseDomain: test.example.com
    privateZone:
      tags:
        Name: test-r69hh-int
        kubernetes.io/cluster/test-r69hh: owned
    publicZone:
      id: Z21IZ5YJJMZ2A4
  status: {}
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""

```

## Monitor for Cluster Completion

```console
$ bin/openshift-install wait-for install-complete
INFO Waiting up to 30m0s for the cluster to initialize...
```

Also, you can observe the running state of your cluster pods:

```console
$ oc get pods --all-namespaces
NAMESPACE                                               NAME                                                                READY     STATUS      RESTARTS   AGE
kube-system                                             etcd-member-ip-10-0-3-111.us-east-2.compute.internal                1/1       Running     0          35m
kube-system                                             etcd-member-ip-10-0-3-239.us-east-2.compute.internal                1/1       Running     0          37m
kube-system                                             etcd-member-ip-10-0-3-24.us-east-2.compute.internal                 1/1       Running     0          35m
openshift-apiserver-operator                            openshift-apiserver-operator-6d6674f4f4-h7t2t                       1/1       Running     1          37m
openshift-apiserver                                     apiserver-fm48r                                                     1/1       Running     0          30m
openshift-apiserver                                     apiserver-fxkvv                                                     1/1       Running     0          29m
openshift-apiserver                                     apiserver-q85nm                                                     1/1       Running     0          29m
...
openshift-service-ca-operator                           openshift-service-ca-operator-66ff6dc6cd-9r257                      1/1       Running     0          37m
openshift-service-ca                                    apiservice-cabundle-injector-695b6bcbc-cl5hm                        1/1       Running     0          35m
openshift-service-ca                                    configmap-cabundle-injector-8498544d7-25qn6                         1/1       Running     0          35m
openshift-service-ca                                    service-serving-cert-signer-6445fc9c6-wqdqn                         1/1       Running     0          35m
openshift-service-catalog-apiserver-operator            openshift-service-catalog-apiserver-operator-549f44668b-b5q2w       1/1       Running     0          32m
openshift-service-catalog-controller-manager-operator   openshift-service-catalog-controller-manager-operator-b78cr2lnm     1/1       Running     0          31m
```

## Add the Ingress DNS Records

If you removed the DNS Zone configuration [earlier](#remove-dns-zones), you'll need to manually create some DNS records pointing at the ingress load balancer.
You can create either a wildcard `*.apps.{baseDomain}.` or specific records (more on the specific records below).
You can use A, CNAME, [alias][route53-alias], etc. records, as you see fit.
For example, you can create wildcard alias records by retrieving the ingress load balancer status:

```console
$ oc -n openshift-ingress get service router-default
NAME             TYPE           CLUSTER-IP      EXTERNAL-IP                                                              PORT(S)                      AGE
router-default   LoadBalancer   172.30.62.215   ab37f072ec51d11e98a7a02ae97362dd-240922428.us-east-2.elb.amazonaws.com   80:31499/TCP,443:30693/TCP   5m
```

Then find the hosted zone ID for the load balancer (or use [this table][route53-zones-for-load-balancers]):

```console
$ aws elb describe-load-balancers | jq -r '.LoadBalancerDescriptions[] | select(.DNSName == "ab37f072ec51d11e98a7a02ae97362dd-240922428.us-east-2.elb.amazonaws.com").CanonicalHostedZoneNameID'
Z3AADJGX6KTTL2
```

And finally, add the alias records to your private and public zones:

```console
$ aws route53 change-resource-record-sets --hosted-zone-id "${YOUR_PRIVATE_ZONE}" --change-batch '{
>   "Changes": [
>     {
>       "Action": "CREATE",
>       "ResourceRecordSet": {
>         "Name": "\\052.apps.your.cluster.domain.example.com",
>         "Type": "A",
>         "AliasTarget":{
>           "HostedZoneId": "Z3AADJGX6KTTL2",
>           "DNSName": "ab37f072ec51d11e98a7a02ae97362dd-240922428.us-east-2.elb.amazonaws.com.",
>           "EvaluateTargetHealth": false
>         }
>       }
>     }
>   ]
> }'
$ aws route53 change-resource-record-sets --hosted-zone-id "${YOUR_PUBLIC_ZONE}" --change-batch '{
>   "Changes": [
>     {
>       "Action": "CREATE",
>       "ResourceRecordSet": {
>         "Name": "\\052.apps.your.cluster.domain.example.com",
>         "Type": "A",
>         "AliasTarget":{
>           "HostedZoneId": "Z3AADJGX6KTTL2",
>           "DNSName": "ab37f072ec51d11e98a7a02ae97362dd-240922428.us-east-2.elb.amazonaws.com.",
>           "EvaluateTargetHealth": false
>         }
>       }
>     }
>   ]
> }'
```

If you prefer to add explicit domains instead of using a wildcard, you can create entries for each of the cluster's current routes:

```console
$ oc get --all-namespaces -o jsonpath='{range .items[*]}{range .status.ingress[*]}{.host}{"\n"}{end}{end}' routes
oauth-openshift.apps.your.cluster.domain.example.com
console-openshift-console.apps.your.cluster.domain.example.com
downloads-openshift-console.apps.your.cluster.domain.example.com
alertmanager-main-openshift-monitoring.apps.your.cluster.domain.example.com
grafana-openshift-monitoring.apps.your.cluster.domain.example.com
prometheus-k8s-openshift-monitoring.apps.your.cluster.domain.example.com
```

[cloudformation]: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/Welcome.html
[delete-stack]: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-console-delete-stack.html
[ingress-operator]: https://github.com/openshift/cluster-ingress-operator
[machine-api-operator]: https://github.com/openshift/machine-api-operator
[route53-alias]: https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/resource-record-sets-choosing-alias-non-alias.html
[route53-zones-for-load-balancers]: https://docs.aws.amazon.com/general/latest/gr/rande.html#elb_region
