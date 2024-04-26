
# Releasing

## Output

### Expected artifacts

1. A container image of the `cluster-api-provider-openstack` controller manager

### Artifact locations

1. The container image is found in the registry `us.gcr.io/k8s-artifacts-prod/capi-openstack/` with an image
   name of `capi-openstack-controller` and a tag that matches the release version. For
   example, in the `v0.2.0` release, the container image location is
   `us.gcr.io/k8s-artifacts-prod/capi-openstack/capi-openstack-controller:v0.2.0`

## Version number

A release version string is: `vX.Y.Z`.

A pre-release version string additionally has a suffix:
- `alpha` for an alpha release
- `beta` for a beta release
- `rc` for a release candidate
and an additional index starting at 0. This takes the form: `vX.Y.Z-[suffix].N`. e.g. the first release candidate prior
to version 1.2.3 would be called `v1.2.3-rc.0`.

It is recommended to create at least one release candidate when bumping `X` or `Y`.

## Release notes

Release notes are user visible information providing details of all relevant changes between releases.

The content of the release notes differs depending on the type of release, specifically:

- Stable releases contain a *full* changelog from the last stable release.
- Pre-releases contain a changelog from the previous pre-release, or the last stable release if there isn't one.

## Process

1. Make sure your repo is clean by git's standards. It is recommended to use a fresh checkout.
1. When bumping `X` or `Y` (but not Z or the pre-release suffix) in the release version you must create a new release branch called `release-X.Y`.
   > NOTE: `upstream` should be the name of the remote pointing to `github.com/kubernetes-sigs/cluster-api-provider-openstack`
    - `git checkout main`
    - `git pull`
    - `git checkout -b release-X.Y`
    - `git push --set-upstream upstream`
1. When bumping `X` or `Y` (but not Z or the pre-release suffix) in the release version, ensure you have added a new
   entry to [metadata.yaml](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/main/metadata.yaml)
   as [described in the CAPI book](https://cluster-api.sigs.k8s.io/clusterctl/provider-contract.html#metadata-yaml), and
   that this has been committed to the release branch prior to release.
1. Make sure you are on the correct release branch: `release-X.Y`
1. Set an environment variable with the version, e.g.:
    - `VERSION=v0.6.0`
1. Create an annotated tag
    - `git tag -s -a $VERSION -m $VERSION`.
1. Push the tag to the GitHub repository:
   > NOTE: `upstream` should be the name of the remote pointing to `github.com/kubernetes-sigs/cluster-api-provider-openstack`
    - `git push upstream $VERSION`

   This will cause the image to be automatically built by CI and pushed to the staging repository. As this only builds
   the image, it only takes a few minutes.
1. Follow the [image promotion process](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/README.md#image-promoter) to promote the image from the staging repo to `registry.k8s.io/capi-openstack`.
   The staging repository can be inspected at https://console.cloud.google.com/gcr/images/k8s-staging-capi-openstack/GLOBAL. Be
   sure to choose the top level `capi-openstack-controller`, which will provide the multi-arch manifest, rather than one for a specific architecture.
   The image build logs are available at [Cloud Build](https://console.cloud.google.com/cloud-build/builds?project=k8s-staging-capi-openstack).
   Add the new sha=>tag mapping to the [images.yaml](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/images.yaml) (use the sha of the image with the corresponding tag). The PR to update the [images.yaml](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/images.yaml) must be approved in the [OWNERS](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/OWNERS) file and merged.

   It is good practise to get somebody else to review this PR. It is safe to perform the following steps while waiting
   for review and the promotion of the image.
1. Run `make release` to build artifacts to be attached to the GitHub release.
1. Generate and finalize the release notes and save them for the next step.
   - Run `make release-notes RELEASE_NOTES_ARGS="--from <tag>"`.
     - Depending on the type of release, substitute `<tag>` with the following:
       - Stable releases: tag of the last stable release
       - Pre-releases*: tag of the latest pre-release (or last stable release if there isn't one)
   - Pay close attention to the `## :question: Sort these by hand` section, as it contains items that need to be manually sorted.
1. Create a draft release in GitHub based on the tag created above
    - Name the release `Release [VERSION]` where VERSION is the full version string.
    - For a pre-release, ensure you check `This is a pre-release` in GitHub when creating the release.
    - Add the release notes generated above as the release description.
1. Attach the following files to the draft release:
    - `./out/infrastructure-components.yaml`
    - `./out/cluster-template.yaml`
    - `./out/cluster-template-without-lb.yaml`
    - `./out/cluster-template-flatcar.yaml`
    - `./out/metadata.yaml`
1. Ensure that the release image has been promoted.
1. Publish release.

### Permissions

Releasing requires a particular set of permissions.

* Approver role for the image promoter process ([kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/OWNERS](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/OWNERS))
* Tag push and release creation rights to the GitHub repository (team `cluster-api-provider-openstack-maintainers` in [kubernetes/org/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml](https://github.com/kubernetes/org/blob/main/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml))

## Staging

There is a post-submit Prow job running after each commit on `main` which pushes a new image to the staging repo (`gcr.io/k8s-staging-capi-openstack/capi-openstack-controller:latest`). Following configuration is involved:
* staging gcr bucket: [kubernetes/k8s.io/blob/main/registry.k8s.io/manifests/k8s-staging-capi-openstack/promoter-manifest.yaml](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/manifests/k8s-staging-capi-openstack/promoter-manifest.yaml)
* post-submit `post-capi-openstack-push-images` Prow job: [kubernetes/test-infra/blob/master/config/jobs/image-pushing/k8s-staging-cluster-api.yaml](https://github.com/kubernetes/test-infra/blob/master/config/jobs/image-pushing/k8s-staging-cluster-api.yaml)) (corresponding dashboard is located at [https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#post-capi-openstack-push-images](https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#post-capi-openstack-push-images))
* Google Cloud Build configuration which is used by the Prow job: [kubernetes-sigs/cluster-api-provider-openstack/cloudbuild.yaml](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/main/cloudbuild.yaml)
