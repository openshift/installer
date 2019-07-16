# Alternative Release-Image Sources

## release-image content vs release-image source

`release-image` content is what operators get installed and what version of each operator (i.e. the container image) gets installed to the cluster. While, `release-image` source is **where** the `release-image` content gets pulled.

Currently the installer controls **both** the `release-image` content and the `release-image` source using the embedded release-image location or the `OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE` env, but to allow users to use a release-image from their private registry (change the `release-image` source), but keep the `release-image` content identical to one vetted by OpenShift, the installer needs to allow specifying source separate from the content.

## Controlling the content

The content of the `release-image`, i.e. the digest, continues to be controlled by the embedded release-image location or the `OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE` env.

## Controlling the source

The installer allows the users to specify the sources for the release-image repository and the other repositories referenced in the release-image through the InstallConfig.

## Details

The design is based on the assumption that all flows of using multiple sources/repositories for the release-image originate from user's using the `oc adm release mirror` command to create those sources/repositories.

## InstallConfig

```go
type InstallConfig struct {
    ///

    // ImageContentSources is the list of sources/repositories that can be used to pull the same content.
    // No two ImageContentSource in the list can include the same repository. Each ImageContentSource must be a disjoint set from the rest.
    ImageContentSources []ImageContentSource `json:"imageContentSources"`

    ///
}

// ImageContentSource defines a list of sources/repositories that can be used to pull a content.
type ImageContentSource struct {
    Sources []string `json:"Sources"`
}
```

If the release-image `q.io/ocp/release-x.y@sha256:abc` which has references to the images in repositories `q.io/openshift/x.y` was mirrored to `local.registry.com/ocp/release-x.y`, the `install-config.yaml` would look like,

```yaml
...
imageContentSources:
- sources:
  - local.registry.com/ocp/release-x.y
  - q.io/ocp/release-x.y
  - q.io/openshift/x.y
...
```

### ImageContentSourcePolicy

If a list of `ImageContentSources` is specified, the installer configures the [RepositoryDigestMirrors][repository-digest-mirrors] for each `ImageContentSource`.

If the `install-config.yaml` provided by the user is:

```yaml
...
imageContentSources:
- sources:
  - local.registry.com/ocp/release-x.y
  - q.io/ocp/release-x.y
  - q.io/openshift/x.y
...
```

The `ImageContentSourcePolicy` object would look like:

```yaml
...
repositoryDigestMirrors:
- sources:
  - local.registry.com/ocp/release-x.y
  - q.io/ocp/release-x.y
  - q.io/openshift/x.y
...
```

### release-image location

The release-image location that is propagated to the bootstrap node and the cluster-version-operator will continue to be the embedded release-image location.

### Bootstrap machine containers-registries.conf

If a list of `ImageContentSources` is specified, the [Registries][registry-containers-registries-conf] will be configured to have each source be a mirror for another.

For example,

If the `install-config.yaml` provided by the user is:

```yaml
...
imageContentSources:
- sources:
  - local.registry.com/ocp/release-x.y
  - q.io/ocp/release-x.y
  - q.io/openshift/x.y
...
```

```toml
[[registry]]
location = "local.registry.com/ocp/release"
mirror-by-digest-only = true

[[registry.mirror]]
location = "local.registry.com/ocp/release"

[[registry.mirror]]
location = "q.io/ocp/release-x.y"

[[registry.mirror]]
location = "q.io/openshift/x.y"

[[registry]]
location = "q.io/ocp/release-x.y"
mirror-by-digest-only = true

[[registry.mirror]]
location = "local.registry.com/ocp/release"

[[registry.mirror]]
location = "q.io/ocp/release-x.y"

[[registry.mirror]]
location = "q.io/openshift/x.y"

[[registry]]
location = "q.io/openshift/x.y"
mirror-by-digest-only = true

[[registry.mirror]]
location = "local.registry.com/ocp/release"

[[registry.mirror]]
location = "q.io/ocp/release-x.y"

[[registry.mirror]]
location = "q.io/openshift/x.y"
```

### oc adm release mirror

The `release mirror` mirrors the a release-image and all other images referenced in the release-image to another repository and then provides user's details of setting up the `install-config.yaml` for the user.

The output would look like:

```console
Release Image q.io/ocp/release-x.y@sha256:abcd was successfully mirrored to local.registry.com/ocp/release-x.y@sha256:abcd

Following section can be added to the install-config.yaml to create a cluster using new repository:
imageContentSources:
- sources:
  - local.registry.com/ocp/release-x.y
  - q.io/ocp/release-x.y
  - q.io/openshift/x.y
```

[repository-digest-mirrors]: https://github.com/openshift/api/blob/de5ca909c7322bb8d06fa5a9e5604491b373da52/operator/v1alpha1/types_image_content_source_policy.go#L50
[registry-containers-registries-conf]: https://github.com/containers/image/blob/v2.0.0/docs/containers-registries.conf.5.md#remapping-and-mirroring-registries
