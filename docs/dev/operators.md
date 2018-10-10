# How to add a new operator to the installer?

This document describes how to provide an operator's configuration files and manifests to installer-launched clusters.

One can classify all the manifest/config files of an operator in two categories:
 - Static
 - Templatized

The static ones are clearly just byte blobs and the templates are ones that have some variable that needs to be filled up. Typically the variables that need filling up will come from user preferences as specified in the install config (e.g. cluster name, cluster domain, service cidr etc). Other dependencies could be TLS cert keys for example.

## The recommended way
The static and template files of your operator need to be dealt with separately.

### Static files

For static files use the Cluster Version Operator (CVO) payload mechanism. There is a particular way to keep the manifest files so that the CVO update payload can pick it up.
See this document:

https://github.com/openshift/cluster-version-operator/tree/master/docs/dev/operators.md

Also remember that the order of creation of the files is alphabetical, so the files should be numbered like below to effectively create the service account before the deployment.
```
00-namespace.yaml
...
03-roles.yaml
04-serviceaccount.yaml
05-deployment.yaml
```
where, 04/05 is the internal ordering of the resource manifests.

### What to do for the dynamic template files?

An operator should auto-discover the install config rather than expand the templates through the installer integration (see alternative). Simply use the apiserver access to get the install-config as a config map. The config map is stored by the name ‘cluster-config-v1’ in the ‘kube-system’ namespaces.
Pseudo code:
```
kube-client(apiserver-url).Get(Resource: "config-map", Namespace: "kube-system", Name: "cluster-config-v1")
```
where, apiserver-url is a cluster supplied ENV var ‘KUBERNETES_SERVICE_HOST’ in the pod.
Example:

https://github.com/openshift/machine-config-operator/blob/e932afdec07dc86d5b643590164f86811e205c57/pkg/operator/operator.go#L272

After discovering the InstallConfig, the operator pod can do two things:

 - create its own configuration in memory as users should not be editing it
 - or, push the discovered config to an API as the operator's users might want to change it in the future

See an example of the configuration for the operator being discovered rather than from a configmap:

https://github.com/openshift/machine-config-operator/blob/31c20eefca172d5c1173e7b79b30bad540958dfe/pkg/operator/render.go#L46

## The alternative (for exceptions only)

Creating a new asset in the installer source is only for exceptional cases where CVO cannot take all the required manifests/config files, and auto-discover is not possible or insufficient e.g. the network operator. As another example, the machine-config-operator needs TLS config.

Such operators need to be directly integrated in the installer's [`manifests` package](../../pkg/asset/manifests). Within this, there are two ways to get the manifests/config files integrated:

 - A new asset for the operator
Create a new operator asset, and render the Dependencies, Name, Load and Generate functions. The Dependencies might contain InstallConfig as an example. In the Generate function, create the config files as asset contents. For the config/manifest’s actual structure, one can choose to vendor the operator’s github pkg, or, if the configuration structures are fairly simple then just copy the definitions directly and avoid the hassle of vendoring. Finally, return the entire list of configs and manifests in the Generate function. 

 - Template files
In the pkg/asset/manifests/content/tectonic directory, place the templates golang variables. Then modify pkg/asset/manifests/tectonic.go to expand the template. Expand templateData in template.go for filling up the template variables.
