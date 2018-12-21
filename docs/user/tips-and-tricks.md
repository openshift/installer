# Tips and Tricks

## Reusing an Install Config

By default, the installer prompts for the necessary information every time a cluster is created. While convenient for one-off cases, this can become tiring if many clusters are needed. The prompts can be bypassed by taking advantage of the installer's notion of [multiple-invocations].

Start by creating an install config and saving it in a cluster-agnostic location:

```console
openshift-install create install-config --dir=initial
mv initial/install-config.yaml .
rm -rf initial
```

Future clusters can then be created by copying that install config into the target directory and then invoking the installer:

```console
mkdir cluster-0
cp install-config.yaml cluster-0/
openshift-install create cluster --dir=cluster-0
```

[multiple-invocations]: overview.md#multiple-invocations
