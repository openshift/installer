# Modifying Tectonic Installer

Modifications are great, but some might result in failed cluster creation, or be incompatible with the Tectonic Installer development process. This document provides an outline of changes to the Terraform modules, configs, and manifests included in Tectonic Installer that can and can not be modified successfully.

Please also note that using "alpha features" through existing Beta or Stable APIs (even on your local resources) is discouraged in production. There is no guarantee that these features will survive a cluster upgrade.

## Machine level modifications

Always safe to modify:
* Container Linux channels. All Container Linux channels may be modified (stable, beta, and alpha).

Never safe to modify:
* Kubelet configuration, including CNI. Modification of the kubelet configuration may result in an inability to start pods, or a failure in communication between cluster components.

May be safe to modify, but must be managed individually:
* Changes to Ignition profiles, such as networking and mounting storage. The process by which these changes within your local fork will be merged back into a new release of Tectonic Installer has not yet been defined.

## Kubernetes level modifications

Always safe to modify:
* Pods, deployments, and any other local component of the cluster. It is always safe to modify your local instance of the Kubernetes cluster.
* Namespaces. It is always safe to add or modify local namespaces.
* Custom RBAC roles. It is always safe to add or modify local RBAC roles.

Never safe to modify:
* Any manifests in `kube-system` or `tectonic-system`. Modifications to these manifests may result in an inability to perform a cluster upgrade.
* Default RBAC roles. Modifications to the default RBAC roles may prevent cluster control plane components from functioning.

## Infrastructure level modifications

Always safe to modify:

Never safe to modify:
* Security group settings.
* Role permissions. Cloud Provider role permissions must meet and exceed documented requirements. (For example: [AWS IAM][iam] or [Azure RBAC][rbac].)
* EC2 block device mapping.
* EC2 AMIs.

Modifying any of these settings might lead to invalid clusters.


[iam]: http://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html
[rbac]: https://docs.microsoft.com/en-us/azure/active-directory/role-based-access-control-what-is
