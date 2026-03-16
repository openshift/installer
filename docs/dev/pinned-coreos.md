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

`data/data/coreos/rhcos.json` is a multi-stream JSON file keyed by stream name
(e.g. `"rhcos-4.21"`, `"rhel-coreos-10"`).  To update a stream, extract it,
update it with `plume`, and merge it back:

```bash
# Step 1: Extract the stream to update into a temp file
jq '.["rhcos-4.21"]' data/data/coreos/rhcos.json > /tmp/stream-update.json

# Step 2: Update the extracted stream with plume (unchanged plume invocation)
plume cosa2stream \
  --target /tmp/stream-update.json \
  --distro rhcos \
  x86_64=9.6.YYYYMMDD-0 aarch64=9.6.YYYYMMDD-0 \
  --url https://rhcos.mirror.openshift.com/art/storage/prod/streams/rhel-9.6 \
  --no-signatures

# Step 3: Merge updated stream back into rhcos.json
tmp=$(mktemp)
jq --slurpfile s /tmp/stream-update.json '.["rhcos-4.21"] = $s[0]' \
  data/data/coreos/rhcos.json > "$tmp" && mv "$tmp" data/data/coreos/rhcos.json
```

For adding or updating the RHCOS 10 stream (run by ART with real build versions):

```bash
# Step 1: Generate stream metadata into a temp file
plume cosa2stream \
  --name rhel-coreos-10 \
  --target /tmp/rhcos10.json \
  --distro rhcos \
  x86_64=10.0.YYYYMMDD-0 aarch64=10.0.YYYYMMDD-0 \
  --url https://rhcos.mirror.openshift.com/art/storage/prod/streams/rhel-10.0 \
  --no-signatures

# Step 2: Merge into rhcos.json
tmp=$(mktemp)
jq --slurpfile s /tmp/rhcos10.json '.["rhel-coreos-10"] = $s[0]' \
  data/data/coreos/rhcos.json > "$tmp" && mv "$tmp" data/data/coreos/rhcos.json
```

For more information on `plume cosa2stream`, see:

- https://github.com/coreos/coreos-assembler/pull/2000
- https://github.com/coreos/coreos-assembler/pull/2052
### Origin of stream metadata


For historical reference, the initial file `data/data/rhcos-stream.json` was generated this way:

```
$ plume cosa2stream --name rhcos-4.8 --distro rhcos  x86_64=48.83.202102230316-0 s390x=47.83.202102090311-0 ppc64le=47.83.202102091015-0 > data/data/rhcos-stream.json
```
NOTE: the data for `data/data/rhcos-stream.json` now lives in `data/data/coreos/rhcos.json`
