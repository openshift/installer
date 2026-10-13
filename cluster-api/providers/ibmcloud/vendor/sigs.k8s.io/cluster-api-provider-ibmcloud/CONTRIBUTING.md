<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Contributing guidelines](#contributing-guidelines)
  - [Branches](#branches)
    - [Support and guarantees](#support-and-guarantees)
  - [Sign the CLA](#sign-the-cla)
  - [Contributing A Patch](#contributing-a-patch)
  - [Issue and Pull Request Management](#issue-and-pull-request-management)
  - [Pre-check before submitting a PR](#pre-check-before-submitting-a-pr)
  - [Build and push images](#build-and-push-images)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Contributing guidelines

## Branches

Cluster API Provider IBM Cloud has two types of branches: the *main* branch and
*release-X* branches.

The *main* branch is where development happens. All the latest and
greatest code, including breaking changes, happens on main.

The *release-X* branches contain stable, backwards compatible code. On every
major or minor release, a new branch is created. It is from these
branches that minor and patch releases are tagged. In some cases, it may
be necessary to open PRs for bugfixes directly against stable branches, but
this should generally not be the case.

### Support and guarantees

Cluster API Provider IBM Cloud maintains the most recent release/releases for all supported API and contract versions. Support for this section refers to the ability to backport and release patch versions;
standard [backport policy](https://github.com/kubernetes-sigs/cluster-api/blob/main/CONTRIBUTING.md#backporting-a-patch) is defined here.

- The API version is determined from the GroupVersion defined in the top-level `api/` package.
- The EOL date of each API Version is determined from the last release available once a new API version is published.

| API Version  | Supported Until      |
|--------------|----------------------|
| **v1beta2**  | TBD (current stable) |
| **v1beta1**  | EOL since 2023-02-09 |

- For the current stable API version (v1beta2) we support the two most recent minor releases; older minor releases are immediately unsupported when a new major/minor release is available.
- For older API versions we only support the most recent minor release until the API version reaches EOL.
- We will maintain test coverage for all supported minor releases for the current stable API version in case we have to do an emergency patch release.
  For example, if v0.5 and v0.6 are currently supported. When v0.7 is released, tests for v0.5 will be removed.

| Minor Release | API Version | Supported Until                                    |
|---------------|-------------|----------------------------------------------------|
| v0.6.x        | **v1beta2** | when v0.8.0 will be released                       |
| v0.5.x        | **v1beta2** | when v0.7.0 will be released, tentatively Nov 2023 |
| v0.4.x        | **v1beta2** | EOL since 2023-09-07 - v0.6.0 release date         |
| v0.3.x        | **v1beta1** | EOL since 2023-02-09 - API version EOL             |

- The CAPI, k8s and test packages will receive regular updates for supported releases to ensure they remain synchronized with the CAPI release being utilized as an integral component of the provider release. This activity is ideally scheduled to occur with every new n-1 and n-2 CAPI minor releases.
- The IBM packages will be monitored for latest updates in conjunction with CAPI minor release update activity, as long as there are no disruptive changes that impact the project stability.
- Exceptions can be filed with maintainers and taken into consideration on a case-by-case basis.

## Sign the CLA

Kubernetes projects require that you sign a Contributor License Agreement (CLA) before we can accept your pull requests.  Please see https://git.k8s.io/community/CLA.md for more info

## Contributing A Patch

1. Submit an issue describing your proposed change to the repo in question.
1. The [repo owners](OWNERS) will respond to your issue promptly.
1. If your proposed change is accepted, and you haven't already done so, sign a Contributor License Agreement (see details above).
1. Fork the desired repo, develop and test your code changes.
1. Submit a pull request.

## Issue and Pull Request Management

Anyone may comment on issues and submit reviews for pull requests. However, in
order to be assigned an issue or pull request, you must be a member of the
[Kubernetes SIGs](https://github.com/kubernetes-sigs) GitHub organization.

If you are a Kubernetes GitHub organization member, you are eligible for
membership in the Kubernetes SIGs GitHub organization and can request
membership by [opening an issue](https://github.com/kubernetes/org/issues/new?template=membership.md&title=REQUEST%3A%20New%20membership%20for%20%3Cyour-GH-handle%3E)
against the kubernetes/org repo.

However, if you are a member of any of the related Kubernetes GitHub
organizations but not of the Kubernetes org, you will need explicit sponsorship
for your membership request. You can read more about Kubernetes membership and
sponsorship [here](https://github.com/kubernetes/community/blob/master/community-membership.md).

Cluster API maintainers can assign you an issue or pull request by leaving a
`/assign <your Github ID>` comment on the issue or pull request.

## Pre-check before submitting a PR

After your PR is ready to commit, please run following commands to check your code:

```shell
make check
make test
```

## Build and push images
Make sure your code build passed:

```shell
export REGISTRY=<your-docker-registry>
make build-push-images
```

Now, you can follow the [getting started guide](./README.md#getting-started) to work with the Cluster API Provider IBM Cloud.
