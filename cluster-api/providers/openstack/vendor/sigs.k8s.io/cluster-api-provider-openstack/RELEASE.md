
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

There is an [issue template](.github/ISSUE_TEMPLATE/new_release.md) to help track release activities.

1. Make sure your repo is clean by git's standards. It is recommended to use a fresh checkout.
1. Repository Setup
   - Clone the repository: `git clone git@github.com:kubernetes-sigs/cluster-api-provider-openstack.git`
   or if using existing repository, make sure origin is set to the fork and
   upstream is set to `kubernetes-sigs`. Verify if your remote is set properly or not
   by using following command `git remote -v`, where origin points to fork and upstream points to main repo.
   - Fetch the remote (`kubernetes-sigs`): `git fetch upstream`
   This makes sure that all the tags are accessible.

1. When bumping `X` or `Y` (but not Z or the pre-release suffix) in the release version, ensure you have added a new
   entry to [metadata.yaml](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/main/metadata.yaml)
   as [described in the CAPI book](https://cluster-api.sigs.k8s.io/developer/providers/contracts/clusterctl#metadata-yaml), and
   that this has been committed to the release branch prior to release.

1. Creating Release Notes
   - Switch to the main branch: `git checkout main`
   - Create a new branch for the release notes**:
     `git checkout -b release-notes-X.Y.Z origin/main`
   - Generate the release notes: `RELEASE_TAG=vX.Y.Z make generate-release-notes`
   - Replace `vX.Y.Z` with the new release tag you're creating.
   - This command generates the release notes here
     `releasenotes/<RELEASE_TAG>.md` .

1. Next step is to clean up the release note manually.
   - If release is not an alpha or a beta or release candidate, check for duplicates,
     reverts, and incorrect classifications of PRs, and whatever release
     creation tagged to be manually checked.
   - For any superseded PRs (like same dependency uplifted multiple times, or
     commit revertion) that provide no value to the release, move them to
     Superseded section. This way the changes are acknowledged to be part of the
     release, but not overwhelming the important changes contained in the release.
   - Commit your changes, push the new branch and create a pull request:
     - The commit and PR title should be 🚀 Release v1.x.y:
       -`git commit -S -s -m ":rocket: Release vX.Y.Z"`
       -`git push -u origin release-notes-X.Y.Z`
     - Important! The commit should only contain the release notes file, nothing
      else, otherwise automation will not work. Push as normal, through your fork (`origin`).
   - Ask maintainers and release team members to review your pull request.

   Once the PR is merged, the following GitHub actions are triggered:

   - GitHub action `Create Release` runs following jobs
     - GitHub job `push_release_tags` will create and push the tags. This action
     will also create release branch if its missing and release is `rc` or minor.
     - GitHub job `create draft release` creates draft release. Don't publish the
     release yet. Running actions are visible on the
      [Actions](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/actions)
      page, and draft release will be visible on top of the
      [Releases](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/releases).

      The image will also be automatically built by CI and pushed to the staging repository. As this only builds the image, it only takes a few minutes.

1. Follow the [image promotion process](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/README.md#image-promoter)
   to promote the image from the staging repo to `registry.k8s.io/capi-openstack`.
   The staging repository can be inspected at [Staging CAPI Openstack](https://console.cloud.google.com/gcr/images/k8s-staging-capi-openstack/GLOBAL). Be
   sure to choose the top level `capi-openstack-controller`, which will provide the multi-arch manifest, rather than one for a specific architecture.
   The image build logs are available at [Cloud Build](https://console.cloud.google.com/cloud-build/builds?project=k8s-staging-capi-openstack).
   Add the new sha=>tag mapping to the [images.yaml](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/images.yaml) (use the sha of the image with the corresponding tag). The PR to update the [images.yaml](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/images.yaml) must be approved in the [OWNERS](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/OWNERS) file and merged.
   Here is an example [pull request](https://github.com/kubernetes/k8s.io/pull/8807).

   It is good practise to get somebody else to review this PR. It is safe to perform the following steps while waiting for review and the promotion of the image.

1. Check carefully the [draft release](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/releases)
   created by the workflow. Ensure that the release notes are correct and that the artifacts are present.
   If any changes are needed, edit the release notes in the GitHub UI and add any missing artifacts.
1. Ensure that the release image has been promoted.
1. If the release you're making is not a new major release, new minor release,
   or a new patch release from the latest release branch, uncheck the box for
   latest release. If it is a release candidate (RC) or a beta or an alpha
   release, tick pre-release box.
1. Publish release.

### Post release actions

1. When bumping `X` or `Y` (but not Z or the pre-release suffix), update the [periodic jobs](https://github.com/kubernetes/test-infra/tree/master/config/jobs/kubernetes-sigs/cluster-api-provider-openstack).
   Make sure there are periodic jobs for the new release branch, and clean up jobs for branches that are no longer supported.
1. When bumping `X` or `Y` (but not Z or the pre-release suffix), update the [clusterctl upgrade tests](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/main/test/e2e/suites/e2e/clusterctl_upgrade_test.go) and the [e2e config](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/main/test/e2e/data/e2e_conf.yaml)
   to include the new release branch.
   It is also a good idea to update the Cluster API versions we test against and to clean up older versions that we no longer want to test.

### Permissions

Releasing requires a particular set of permissions.

1. Approver role for the image promoter process ([kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/OWNERS](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-openstack/OWNERS))

1. Tag push and release creation rights to the GitHub repository (team `cluster-api-provider-openstack-maintainers` in [kubernetes/org/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml](https://github.com/kubernetes/org/blob/main/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml))

## Staging

There is a post-submit Prow job running after each commit on `main` which pushes a new image to the staging repo (`gcr.io/k8s-staging-capi-openstack/capi-openstack-controller:latest`). Following configuration is involved:

1. staging gcr bucket: [kubernetes/k8s.io/blob/main/registry.k8s.io/manifests/k8s-staging-capi-openstack/promoter-manifest.yaml](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/manifests/k8s-staging-capi-openstack/promoter-manifest.yaml)

1. post-submit `post-capi-openstack-push-images` Prow job: [kubernetes/test-infra/blob/master/config/jobs/image-pushing/k8s-staging-cluster-api.yaml](https://github.com/kubernetes/test-infra/blob/master/config/jobs/image-pushing/k8s-staging-cluster-api.yaml) (corresponding dashboard is located at [https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#post-capi-openstack-push-images](https://testgrid.k8s.io/sig-cluster-lifecycle-image-pushes#post-capi-openstack-push-images))

1. Google Cloud Build configuration which is used by the Prow job: [kubernetes-sigs/cluster-api-provider-openstack/cloudbuild.yaml](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/main/cloudbuild.yaml)
