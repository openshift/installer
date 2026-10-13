# Contributing guidelines

## Sign the CLA

Kubernetes projects require that you sign a Contributor License Agreement (CLA) before we can accept your pull requests. Please see https://git.k8s.io/community/CLA.md for more info

### Contributing A Patch

1. Submit an issue describing your proposed change to the repo in question.
1. The [repo owners](OWNERS) will respond to your issue promptly.
1. If your proposed change is accepted, and you haven't already done so, sign a Contributor License Agreement (see details above).
1. Fork the desired repo, develop and test your code changes.
1. Submit a pull request.
1. All code PR must be labeled with ⚠️ (:warning:, major or breaking changes), ✨ (:sparkles:, feature additions), 🐛 (:bug:, patch and bugfixes), 📖 (:book:, documentation or proposals), or 🌱 (:seedling:, minor or other)

## Branches

Cluster API Provider OpenStack has two types of branches: the *main* branch and
*release-X* branches.

The *main* branch is where development happens. All the latest and
greatest code, including breaking changes, happens on main.

The *release-X* branches contain stable, backwards compatible code. On every
major or minor release, a new branch is created. It is from these
branches that minor and patch releases are tagged. In some cases, it may
be necessary to open PRs for bugfixes directly against stable branches, but
this should generally not be the case.

### Support and guarantees

Cluster API Provider OpenStack maintains the most recent release/releases for all supported API and contract versions. Support for this section refers to the ability to backport and release patch versions.

- The API version is determined from the GroupVersion defined in the top-level `api/` package.
- The EOL date of each API Version is determined from the last release available once a new API version is published.

| API Version  | Supported Until       |
|--------------|-----------------------|
| **v1beta1**  | TBD (current stable)  |

- For the current stable API version (v1beta1) we support the two most recent minor releases; older minor releases are immediately unsupported when a new major/minor release is available.
- For older API versions we only support the most recent minor release until the API version reaches EOL.
- We will maintain test coverage for all supported minor releases and for one additional release for the current stable API version in case we have to do an emergency patch release.
  For example, if v0.11 and v0.12 are currently supported, we will also maintain test coverage for v0.10 for one additional release cycle. When v0.13 is released, tests for v0.10 will be removed.

| Minor Release | API Version  | Supported Until                                |
|---------------|--------------|------------------------------------------------|
| v0.12.x       | **v1beta1**  | when v0.14.0 will be released                  |
| v0.11.x       | **v1beta1**  | when v0.13.0 will be released                  |
| v0.10.x       | **v1beta1**  | EOL since 2025-02-06 - v0.12.0 release date    |
| v0.9.x        | **v1alpha7** | EOL since 2024-10-24 - v0.11.0 release date    |
| v0.8.x        | **v1alpha7** | EOL since 2024-04-17 - v0.10.0 release date    |

- Exceptions can be filed with maintainers and taken into consideration on a case-by-case basis.

### Removal of v1alpha apiVersions

| Minor Release | v1beta1       | v1alpha7   | v1alpha6   | v1alpha5   |
|---------------|---------------|------------|------------|------------|
| v0.12.x       | **supported** | not served |            |            |
| v0.11.x       | **supported** | deprecated | not served |            |
| v0.10.x       | **supported** | supported  | deprecated | not served |
| v0.9.x        |               | supported  | supported  | deprecated |
| v0.8.x        |               | supported  | supported  | deprecated |

Note: Removal of a deprecated APIVersion in Kubernetes [can cause issues with garbage collection by the kube-controller-manager](https://github.com/kubernetes/kubernetes/issues/102641)
This means that some objects which rely on garbage collection for cleanup - e.g. MachineSets and their descendent objects, like Machines and InfrastructureMachines, may not be cleaned up properly if those
objects were created with an APIVersion which is no longer served.
To avoid these issues it's advised to ensure a restart to the kube-controller-manager is done after upgrading to a version of which drops support for an APIVersion.
This can be accomplished with any Kubernetes control-plane rollout, including a Kubernetes version upgrade, or by manually stopping and restarting the kube-controller-manager.

Note: We have introduced experimental APIs separate from the "main" API mentioned here.
They do not follow the support cycle described here.
The goal is to mature them separately so that they can be quickly iterated on, dropped or eventually included in the stable API.
