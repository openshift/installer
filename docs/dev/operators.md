# How to add a new operator to the installer?

This document describes how to provide an operator's configuration files and manifests to installer-launched clusters.

## The recommended way

Most operators should use [the Cluster Version Operator (CVO) payload mechanism][cvo-operators].

## The alternative (for exceptions only)

Creating a new asset in the installer source is only for exceptional cases where CVO cannot take all the required manifests/config files, and auto-discover is not possible or insufficient e.g. the network operator. As another example, the machine-config-operator needs TLS config.

Such operators need to be directly integrated in the installer's [`manifests` package](../../pkg/asset/manifests). Within this, there are two ways to get the manifests/config files integrated:

 - A new asset for the operator
Create a new operator asset, and render the Dependencies, Name, Load and Generate functions. The Dependencies might contain InstallConfig as an example. In the Generate function, create the config files as asset contents. For the config/manifest’s actual structure, one can choose to vendor the operator’s github pkg, or, if the configuration structures are fairly simple then just copy the definitions directly and avoid the hassle of vendoring. Finally, return the entire list of configs and manifests in the Generate function. 

 - Template files
In the pkg/asset/manifests/content/openshift directory, place the templates golang variables. Then modify pkg/asset/manifests/openshift.go to expand the template. Expand templateData in template.go for filling up the template variables.

[cvo-operators]: https://github.com/openshift/enhancements/blob/master/dev-guide/cluster-version-operator/dev/operators.md
