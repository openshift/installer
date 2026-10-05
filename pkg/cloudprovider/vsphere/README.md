# pkg/cloudprovider/vsphere

This package is a local copy of `github.com/openshift/library-go/pkg/cloudprovider/vsphere`
(`config.go`, `ini.go`, `yaml.go`, `config_test.go`). It replaces the previous import of that
package in `pkg/asset/manifests/vsphere/cloudproviderconfig.go`.

## Why this is copied instead of imported

[SPLAT-2867](https://redhat.atlassian.net/browse/SPLAT-2867) backports the vSphere multi-vCenter
Day 2 support (parent epic SPLAT-2836) to OpenShift 4.20. Part of that work landed upstream in
library-go as:

- [library-go#2175](https://github.com/openshift/library-go/pull/2175) — shared vSphere cloud
  config modules used by 3CMO (SPLAT-2651)
- [library-go#2195](https://github.com/openshift/library-go/pull/2195) — `CPIConfig` Node struct
  non-pointer fix (SPLAT-2651)

Pulling those changes in the normal way — bumping the `github.com/openshift/library-go` and
`github.com/openshift/api` go.mod requirements to a commit that includes them — drags in the rest
of what moved upstream between the 4.20-pinned commit and that commit (newer `k8s.io/api`,
`sigs.k8s.io/cluster-api`, `google.golang.org/api`, `sigs.k8s.io/kustomize/api`, the new
`cel.dev/expr` module, etc.). Reconciling that transitive graph against the versions the installer
and its other vendored dependencies require on 4.20 turned into a multi-day dependency-resolution
effort with no clean resolution.

To unblock the backport, `go.mod`/`go.sum`/`vendor/` are kept on the library-go and api versions
already used on 4.20, and only the small, self-contained vSphere cloud-config package is copied
into the installer tree here instead of being vendored.

## Keeping this in sync

These files should track the upstream library-go source as closely as possible (see the
"Original code was taken from ..." comments in `ini.go`/`yaml.go`, which describe the further
upstream vsphere-cloud-provider origin of the type definitions). If library-go/api are ever bumped
past the commits containing PRs #2175/#2195 without dependency conflicts, this local copy should
be deleted and `pkg/asset/manifests/vsphere/cloudproviderconfig.go` switched back to importing
`github.com/openshift/library-go/pkg/cloudprovider/vsphere` directly.
