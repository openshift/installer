# Contributing guidelines

## Contributor License Agreements

We'd love to accept your patches! Before we can take them, we have to jump a couple of legal hurdles.

Please fill out either the individual or corporate Contributor License Agreement (CLA). More information about the CLA
and instructions for signing it [can be found here](https://git.k8s.io/community/CLA.md).

***NOTE***: Only original source code from you and other people that have signed the CLA can be accepted into the
repository.

## Versioning

### Branches

CAPV has two types of branches: the *main* branch and *release-X* branches.

The *main* branch is where development happens. All the latest and
greatest code, including breaking changes, happens on main.

The *release-X* branches contain stable, backwards compatible code. On every
major or minor release, a new branch is created. It is from these
branches that minor and patch releases are tagged. In some cases, it may
be necessary to open PRs for bugfixes directly against stable branches, but
this should generally not be the case.

### Backporting a patch

We generally do not accept PRs directly against release branches, while we might accept backports of fixes/changes already
merged into the main branch. In most cases the cherry-pick bot can and should be used to automate opening a cherry-pick PR.

We generally allow backports of following changes to all supported branches:
- Bug fixes and security fixes
- Dependency bumps for CVEs (usually limited to CVE resolution; backports of non-CVE related version bumps are considered exceptions to be evaluated case by case)
- Changes required to support new Kubernetes versions, when possible.
- Changes to use the latest Go patch release. If the Go minor version of a supported branch goes out of support, we will consider on a case-by-case basis
  to bump to a newer Go minor version (e.g. to pick up CVE fixes).
- Improvements to test and CI signal

In addition to that we allow backports at maintainers discretion. Please let us know if you would like us to consider backporting a specific PR.

In general, we support the two latest release branches. In addition, we will keep the CI coverage for older branches around so we're able to cut additional patch 
releases to fix CVEs and critical bugs if needed.

## Contributing A Patch

1. Submit an issue describing your proposed change to the repo in question.
1. The [repo owners](OWNERS) will respond to your issue promptly.
1. If your proposed change is accepted, and you haven't already done so, sign a Contributor License Agreement (see details above).
1. Fork the desired repo, develop and test your code changes.
1. Submit a pull request.
    * All code PR must be labeled with one of
        * ⚠️ (:warning:, major or breaking changes)
        * ✨ (:sparkles:, feature additions)
        * 🐛 (:bug:, patch and bugfixes)
        * 📖 (:book:, documentation or proposals)
        * 🌱 (:seedling:, minor or other)

## Dependency Licence Management

Cluster API provider vSphere follows the [license policy of the CNCF](https://github.com/cncf/foundation/blob/main/allowed-third-party-license-policy.md). This sets limits on which
licenses dependencies and other artifacts use. For go dependencies only dependencies listed in the `go.mod` are considered dependencies. This is in line with [how dependencies are reviewed in Kubernetes](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/vendor.md#reviewing-and-approving-dependency-changes).

### Contributor Ladder

We broadly follow the requirements from the [Kubernetes Community Membership](https://github.com/kubernetes/community/blob/master/community-membership.md).

> When making changes to **OWNER_ALIASES** please check that the **sig-cluster-lifecycle-leads**, **cluster-api-admins** and **cluster-api-maintainers** are correct.

#### Becoming a reviewer

If you would like to become a reviewer, then please ask one of the current maintainers.

We generally try to follow the [requirements for a reviewer](https://github.com/kubernetes/community/blob/master/community-membership.md#reviewer) from upstream Kubernetes. But if you feel that you don't fully meet the requirements then reach out to us, they are not set in stone.

A reviewer can get PRs automatically assigned for review, and can `/lgtm` PRs.

To become a reviewer, ensure you are a member of the **kubernetes-sigs** Github organisation
following [kubernetes/org/issues/new/choose](https://github.com/kubernetes/org/issues/new/choose).

The steps to add someone as a reviewer are:

* Add the GitHub alias to the **cluster-api-vsphere-reviewers** section of [OWNERS_ALIASES](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/main/OWNERS_ALIASES)
* Create a PR with the change that is held (i.e. by using `/hold`)
* Announce the change within the CAPV slack channel and as a PSA in the next CAPV office hours
* After 7 days of lazy consensus or after the next CAPV office hours (whichever is longer) the PR can be merged

#### Becoming a maintainer

If you have made significant contributions to Cluster API Provider vSphere, a maintainer may nominate you to become a maintainer for the project.

We generally follow the [requirements for a approver](https://github.com/kubernetes/community/blob/master/community-membership.md#approver) from upstream Kubernetes. However, if you don't fully meet the requirements then a quorum of maintainers may still propose you if they feel you will make significant contributions.

Maintainers are able to approve PRs, as well as participate in release processes and have write access to the repo. **As a maintainer you will be expected to run the office hours, especially if no one else wants to**.

Maintainers require membership of the **Kubernetes** Github organisation via
[kubernetes/org/issues/new/choose](https://github.com/kubernetes/org/issues/new/choose)

The steps to add someone as a maintainer are:

* Add the GitHub alias to the **cluster-api-vsphere-maintainers** and remove them from **cluster-api-vsphere-reviewers** sections of [OWNERS_ALIASES](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/main/OWNERS_ALIASES)
* Create a PR with the change that is held (i.e. by using `/hold`)
* Announce the change within the CAPV slack channel and as a PSA in the next CAPV office hours
* After 7 days of lazy consensus or after the next CAPV office hours (whichever is longer) the PR can be merged
* Open PR to add Github username to **cluster-api-provider-vsphere-maintainers**
to [kubernetes/org/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml](https://github.com/kubernetes/org/blob/main/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml)
* Open PR to add Github username to [kubernetes/test-infra/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/OWNERS](https://github.com/kubernetes/test-infra/blob/master/config/jobs/kubernetes-sigs/cluster-api-provider-vsphere/OWNERS)
* Open PR to add Google ID to the k8s-infra-staging-capi-vsphere@kubernetes.io and sig-cluster-lifecycle-cluster-api-vsphere-alerts@kubernetes.io Google groups in [kubernetes/k8s.io/groups/sig-cluster-lifecycle/groups.yaml](https://github.com/kubernetes/k8s.io/blob/main/groups/sig-cluster-lifecycle/groups.yaml)
* Open PR to add approvers/reviewers to [CAPV image promotion](https://github.com/kubernetes/k8s.io/blob/main/registry.k8s.io/images/k8s-staging-capi-vsphere/OWNERS).
* Open PR to image-builder to modify `cluster-api-vsphere-maintainers` in [OWNERS_ALIASES](https://github.com/kubernetes-sigs/image-builder/blob/main/OWNERS_ALIASES)

#### Becoming a admin

After a period of time one of the existing CAPV or CAPI admins may propose you to become an admin of the CAPV project.

Admins have GitHub **admin** access to perform tasks on the repo.

The steps to add someone as an admin are:

* Open PR to add Github username to **cluster-api-provider-vsphere-admins**
to [kubernetes/org/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml](https://github.com/kubernetes/org/blob/main/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml)
