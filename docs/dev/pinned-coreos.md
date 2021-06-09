# CoreOS and the installer

A key decision made before the release of OpenShift 4 is to pin the CoreOS bootimage in
the installer: https://github.com/openshift/installer/commit/e080f0494708b2674fe37af02f670c8030c32bf6

That is still the case today; when one gets an `openshift-install` binary, that
binary contains the 2-tuple `(CoreOS, release image)`, meaning the result of an
install will be the same thing each time.

More background:

 - https://github.com/openshift/enhancements/pull/201
 - https://github.com/openshift/machine-config-operator/blob/master/docs/OSUpgrades.md

## Stream metadata

As of 4.8 the [stream metadata enhancement](https://github.com/openshift/enhancements/blob/master/enhancements/coreos-bootimages.md)
is in progress which provides a standardized JSON format and injects
that data into the cluster as well.

### Updating pinned stream metadata


To update the bootimage for one or more architectures, use e.g.

```
$ plume cosa2stream --target data/data/rhcos-stream.json --distro rhcos  x86_64=48.83.202102230316-0 s390x=47.83.202102090311-0 ppc64le=47.83.202102091015-0 --url https://rhcos-redirector.apps.art.xq1c.p1.openshiftapps.com/art/storage/releases
```

For more information on this command, see:

- https://github.com/coreos/coreos-assembler/pull/2000 
- https://github.com/coreos/coreos-assembler/pull/2052

### Updating pinned legacy metadata

To update the legacy metadata, use:

```
./hack/update-rhcos-bootimage.py https://rhcos-redirector.apps.art.xq1c.p1.openshiftapps.com/art/storage/releases/rhcos-4.6/46.82.202008260918-0/x86_64/meta.json amd64
```

This will hopefully be removed soon.

### Origin of stream metadata


For historical reference, the initial file `data/data/rhcos-stream.json` was generated this way:

```
$ plume cosa2stream --name rhcos-4.8 --distro rhcos  x86_64=48.83.202102230316-0 s390x=47.83.202102090311-0 ppc64le=47.83.202102091015-0 > data/data/rhcos-stream.json
```
