# Tectonic Platform SDK

![Lifecycle Prototype](https://img.shields.io/badge/Lifecycle-Prototype-f4cccc.svg)

The Tectonic Platform SDK provides pre-built recipes to help users create the underlying compute infrastructure for a [Self-Hosted Kubernetes Cluster](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/self-hosted-kubernetes.md) ([vid](https://coreos.com/blog/self-hosted-kubernetes.html)) using [Hashicorp Terraform](https://terraform.io), [bootkube](https://github.com/kubernetes-incubator/bootkube), and supporting tooling.

The goal is to provide well-tested defaults that can be customized for various environments and plugged into other systems.

The unique power of Self-Hosted Kubernetes is that it cleanly separates out the infrastructure from Kubernetes enabling this separation of concerns:

![](http://i.imgur.com/Gd9W7RR.gif)

## Getting Started

Generally:

1. Use the tectonic installer to configure an AWS cluster.
2. Go through the process, do not apply the configuration, but download the assets manually.
3. Unzip the assets in this directory:

```
$ unzip ~/Downloads/<name>-assets.zip
```

## Openstack

### Nova network

Prerequsities:

1. Since openstack nova doesn't provide any DNS registration service, AWS Route53 is being used.
Ensure you have a configured `aws` CLI installation.
2. Ensure you have OpenStack credentials set up, i.e. the environment variables `OS_TENANT_NAME`, `OS_USERNAME`, `OS_PASSWORD`, `OS_AUTH_URL`, `OS_REGION_NAME` are set.

Modify the unzipped assets and substitute `<name>` with the chosen tectonic cluster name and `<base_domain>` with the chosen tectonic DNS base domain:

1. Remove the `--cloud-provider=aws` setting from the apiserver and controller-manager manifests.
2. Replace the ingress URL in all tectonic assets of the cluster to use the nginx ingress node port `32000` since nova doesn't bring a load balancer.

```
diff -r /home/user/original/assets/manifests/kube-apiserver.yaml assets/manifests/kube-apiserver.yaml
42d41
<         - --cloud-provider=aws
44c43
<         - --oidc-issuer-url=https://<name>.<base_domain>/identity
---
>         - --oidc-issuer-url=https://<name>.<base_domain>:32000/identity
diff -r /home/user/original/assets/manifests/kube-controller-manager.yaml assets/manifests/kube-controller-manager.yaml
29d28
<         - --cloud-provider=aws
diff -r /home/user/original/assets/tectonic/console-deployment.yaml assets/tectonic/console-deployment.yaml
41c41
<           value: https://<name>.<base_domain>
---
>           value: https://<name>.<base_domain>:32000
49c49
<           value: https://<name>.<base_domain>/identity
---
>           value: https://<name>.<base_domain>:32000/identity
diff -r /home/user/original/assets/tectonic/identity-config.yaml assets/tectonic/identity-config.yaml
8c8
<     issuer: https://<name>.<base_domain>/identity
---
>     issuer: https://<name>.<base_domain>:32000/identity
23c23
<       - 'https://<name>.<base_domain>/auth/callback'
---
>       - 'https://<name>.<base_domain>:32000/auth/callback'
```

Invoke:

```
$ terraform apply -var 'cluster_name=<name>' openstack-novanet
...
null_resource.copy_assets: Still creating... (5m50s elapsed)
null_resource.copy_assets: Creation complete

Apply complete! Resources: 43 added, 0 changed, 0 destroyed.
```

The tectonic cluster will be reachable under `https://<name>.<base_domain>:32000`.

### Neutron network

The instructions for Neutron are the same as for Nova.

**TODO(sur): verify**

## AWS

*Common Prerequsities*
* Configure AWS credentials via environment variables. 
[See docs](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-environment)
* Configure a region by setting `AWS_REGION` environment variable
* Run through the official Tectonic intaller steps without clicking `Submit` on the last step. Instead click on `Manual boot` below to download the assets zip file.
* Populate a cluster configuration in `config.tfvars` following `example-config.tfvars`

### Using Autoscaling groups

1. Ensure all *prerequsities* are met.
2. From the root of the repo, run `make apply PLATFORM=aws-asg ASSETS=<name-of-assets.zip>`

To clean up run `make destroy PLATFORM=aws-asg ASSETS=<name-of-assets.zip>`

### Without Autoscaling groups

1. Ensure all *prerequsities* are met.
2. From the root of the repo, run `make apply PLATFORM=aws-noasg ASSETS=<name-of-assets.zip>`

To clean up run `make destroy PLATFORM=aws-noasg ASSETS=<name-of-assets.zip>`

**TODO(alexsomesan)**

## Roadmap

This is an unprioritized list of future items the team would like to tackle:

- Run [Kubernetes e2e tests](https://github.com/coreos-inc/tectonic-platform-sdk/issues/6) over repo automatically
- Build a tool to walk the Terraform graph and warn if cluster won't comply with [Generic Platform](https://github.com/coreos-inc/tectonic-platform-sdk/blob/master/Documentation/generic-platform.md)
- Additional platforms like Azure, VMware, Google Cloud, etc
- Create a spec for generic and platform specific Terraform Variable files
- Document how to customize each of the platforms
- Create a tool to verify Terraform Variable files
- Deploy with other self-hosted tools like kubeadm
- Terraform plugin and integration with [matchbox](https://github.com/coreos/matchbox) for bare metal
