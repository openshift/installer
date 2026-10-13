# Alternative Release-Image Sources

## release-image content vs release-image source

`release-image` **content** - which operators get installed to the cluster & which version/container image of each operator  
`release-image` **source** - where `release-image` content gets pulled from

The installer controls *both* the `release-image` content and source using the embedded release-image location or the `OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE` env. Users may also use a release-image from a private registry (change the `release-image` source), but keep the `release-image` content identical to one vetted by OpenShift.

## Controlling the content

The content of the `release-image`, i.e. the digest, continues to be controlled by the embedded release-image location or the `OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE` env.

## Controlling the source

The installer allows the users to specify sources for the release-image repository and other repositories referenced in the release-image through the InstallConfig.

## Details

The design is based on the assumption that all flows of using multiple sources/repositories for the release-image originate with the `oc adm release mirror` command to create those sources/repositories.

## InstallConfig

```go
type InstallConfig struct {
    // ImageContentSources lists sources/repositories for the release-image content.
    ImageContentSources []ImageContentSource `json:"imageContentSources"`
}

// ImageContentSource defines a list of sources/repositories that can be used to pull content.
type ImageContentSource struct {
    Source  string   `json:"source"`
    Mirrors []string `json:"mirrors"`
}
```

If the release-image `q.io/ocp/release-x.y@sha256:abc` which has references to the images in repositories `q.io/openshift/x.y` was mirrored to `local.registry.com/ocp/release-x.y`, the `install-config.yaml` would look like,

```yaml
...
imageContentSources:
- source: q.io/ocp/release-x.y
  mirrors:
  - local.registry.com/ocp/release-x.y
- source: q.io/openshift/x.y
  mirrors:
  - local.registry.com/ocp/release-x.y
...
```

### ImageContentSourcePolicy

If a list of `ImageContentSources` is specified, the installer configures the [RepositoryDigestMirrors][repository-digest-mirrors] for each `ImageContentSource`.

Using the same `install-config.yaml` from above, the `ImageContentSourcePolicy` object would look like:

```yaml
...
repositoryDigestMirrors:
- source: q.io/ocp/release-x.y
  mirrors:
  - local.registry.com/ocp/release-x.y
- source: q.io/openshift/x.y
  mirrors:
  - local.registry.com/ocp/release-x.y
...
```

### release-image location

The release-image location that is propagated to the bootstrap node and the cluster-version-operator will continue to be the embedded release-image location.

### Bootstrap machine containers-registries.conf

If a list of `ImageContentSources` is specified, the [Registries][registry-containers-registries-conf] will be configured to have each source be a mirror for another.

For example, our `install-config.yaml` will result in:

```toml
[[registry]]
location = "q.io/ocp/release-x.y"
mirror-by-digest-only = true

[[registry.mirror]]
location = "local.registry.com/ocp/release"

[[registry]]
location = "q.io/openshift/x.y"
mirror-by-digest-only = true

[[registry.mirror]]
location = "local.registry.com/ocp/release"
```

### oc adm release mirror

The `release mirror` mirrors the release-image and all other images referenced in the release-image to another repository and then provides details for setting up the `install-config.yaml`.

The output would look like:

```console
Release Image q.io/ocp/release-x.y@sha256:abcd was successfully mirrored to local.registry.com/ocp/release-x.y@sha256:abcd

Following section can be added to the install-config.yaml to create a cluster using new repository:
imageContentSources:
- source: q.io/ocp/release-x.y
  mirrors:
  - local.registry.com/ocp/release-x.y
- source: q.io/openshift/x.y
  mirrors:
  - local.registry.com/ocp/release-x.y
```

[repository-digest-mirrors]: https://github.com/openshift/api/blob/9525304a0adb725ab4a4a54539a1a6bf6cc343d3/operator/v1alpha1/types_image_content_source_policy.go#L56
[registry-containers-registries-conf]: https://github.com/containers/image/blob/v2.0.0/docs/containers-registries.conf.5.md#remapping-and-mirroring-registries
