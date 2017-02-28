# tectonic-platform-sdk

Tectonic Platform SDK (Will be Open Sourced in the Future)

Documentation is evolving quickly on this Google Doc: https://docs.google.com/document/d/1msWuuMsIfZMvzs2qZbOLOKEwciJFKh2iFImLxnTmyVg/edit#heading=h.vuc6f8srgdo

# Installation

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

### Using Autoscaling groups

**TODO(alexsomesan)**

### Without Autoscaling groups

**TODO(alexsomesan)**
