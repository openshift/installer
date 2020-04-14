## IAM Instance Role Permissions

Historically, the installer has created a set of IAM instance permissions. 
This list was not managed by the cluster, was not updated for any new use cases and can be considered permissive due
to the usage of wildcards and automatically inheriting new capabilities as they were added by AWS. 

Installations of OpenShift will now receive a tightened set of permissions matching the current in-tree provider.
It is backwards compatible with all versions of OpenShift to use the new set of permissions.
Additional AWS cloud capabilities and IAM permissions will not be enabled or added until IAM instance role permissions
come under management of an operator or can be eliminated entirely by the use of discrete pod identity directly 
assigned to running services within OpenShift. 

For all other uses/needs of IAM credentials, please see the [cloud-credential-operator](#cco).

[cco]: https://github.com/openshift/cloud-credential-operator