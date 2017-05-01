# Uninstall

## Delete LoadBalancer services and CloudFormation StackFormation

### Delete Kubernetes LoadBalancer services

An attempt to delete a CloudFormation stack on which a Tectonic cluster has been deployed may hang if there are existing Kubernetes `LoadBalancer` services in the cluster.

Kubernetes [`LoadBalancer` services][k8s-lb] expose cluster facilities by manipulating the AWS Elastic Load Balancer (ELB) API. These ELBs are created dynamically, after the creation of the CloudFormation stack that defines the cluster nodes and network environment. Such ELBs that are not part of the CloudFormation stack, but that are attached to a CloudFormation-managed subnet, can cause attempts to delete the CloudFormation stack to fail.

Delete all services of `type=LoadBalancer` with `kubectl delete` or in Tectonic Console to avoid this issue.

For example, to delete *every* LoadBalancer service in a cluster with `kubectl`:

```sh
$ kubectl delete services -l type=LoadBalancer
```

### Delete other resources attached to but not part of the CloudFormation stack

Kubernetes LoadBalancer service ELBs are the most common resource added after Tectonic installation, but other resources that are attached to the Tectonic CloudFormation stack, but not defined in it, can also interfere with deleting the stack.

Such resources include ELBs manually created with endpoints in the CloudFormation subnet, Elastic Block Store (EBS) volumes that were attached to nodes or pods after installing Tectonic, and any other resources associated with Tectonic's Virtual Private Cloud (VPC) but not created by the Tectonic installer. Delete these resources by using their AWS service's UI or API.

Attempts to delete a CloudFormation stack that have failed can usually be restarted and successfully completed after deleting any `LoadBalancer` services and other post-install AWS resources attached to the Tectonic VPC.

### Delete the CloudFormation StackFormation

Find the CouldFormation StackFormation ID that was provisioned during the installation process.

[Delete the StackFormation][aws-delete-stackf]. This will take 10-20 minutes.

If any errors occur be sure to delete any additional resources you may have added that reference resources in the StackFormation.

## Reinstall

Use the latest [Tectonic AWS][install-aws] documentation to create a new cluster and deploy your applications and services. If your new cluster will be similar to a previous cluster, you may supply the `tectonic.progress` "progress file" found in your existing [assets bundle][assets] to base the installation on that configuration.


[assets]: ../../admin/assets-zip.md
[aws-delete-stackf]: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-console-delete-stack.html
[install-aws]: index.md
[k8s-lb]: https://kubernetes.io/docs/user-guide/load-balancer/
