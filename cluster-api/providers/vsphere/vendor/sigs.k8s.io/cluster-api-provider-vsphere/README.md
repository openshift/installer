# Kubernetes Cluster API Provider vSphere

[![Go Report Card](https://goreportcard.com/badge/github.com/kubernetes-sigs/cluster-api-provider-vsphere)](https://goreportcard.com/report/github.com/kubernetes-sigs/cluster-api-provider-vsphere)

<img src="https://github.com/kubernetes/kubernetes/raw/master/logo/logo.png" width="100" height="100" /><a href="https://www.vmware.com/products/vsphere.html"><img height="100" hspace="90px" src="https://i.imgur.com/Wd24COX.png" alt="Powered by VMware vSphere" /></a>

Kubernetes-native declarative infrastructure for vSphere.

## What is the Cluster API Provider vSphere

The [Cluster API][cluster_api] brings declarative, Kubernetes-style APIs to cluster creation, configuration and management. Cluster API Provider for vSphere is a concrete implementation of Cluster API
for vSphere.

The API itself is shared across multiple cloud providers allowing for true vSphere hybrid deployments of Kubernetes. It is built atop the lessons learned from previous cluster managers such
as [kops][kops] and [kubicorn][kubicorn].

## Launching a Kubernetes cluster on vSphere

Check out the [getting started guide](./docs/getting_started.md) for launching a cluster on vSphere.

## Features

- Native Kubernetes manifests and API
- Manages the bootstrapping of VMs on cluster.
- Choice of Linux distribution between Ubuntu 18.04 and CentOS 7 using VM Templates based on [OVA images](#kubernetes-versions-with-published-ovas).
- Deploys Kubernetes control planes into provided clusters on vSphere.
- Doesn't use SSH for bootstrapping nodes.
- Installs only the minimal components to bootstrap a control plane and workers.

------

# CAPV release support

A Cluster API provider vSphere release correspond to a release in the GitHub repository for this project, and the corresponding images published in the Kubernetes docker registry.
The CAPV release calendar will be aligned as much as possible to the Kubernetes and the Cluster API release dates plus an offset for 
performing test/validation after K8s/CAPI release are available (tentatively one or two weeks after both K8s and CAPI release are available). 

The CAPV team actively supports the latest two minor releases (N, N-1); support in this context means that we:

- Have CI signal with E2E tests, unit tests, CVE scans etc.
- Accept bug fixes, perform golang or dependency bumps, etc.
- Cut patch releases when necessary.

On top of supporting the N and N-1 releases, the CAPV team also maintains CI signal for the CAPV N-2 releases in case we have to do an emergency patch release. Please note that:

- Even if a subset of the CI signal for the N-2 branch is preserved, the N-2 branch is considered in maintenance mode and no change is back-ported proactively.
- Security scans will be disabled (this signal does not make sense considering that CVE are not going to be fixed proactively)
- If there is a need for an emergency patch, e.g. to fix a critical issue, please bring this up to maintainers and it will be considered on a case-by-case basis.

All considered, each CAPV minor release is supported for a period of roughly 12 months:

- The first eight months of this timeframe will be considered the standard support period for a minor release.
- The next four months the minor release will be considered in maintenance mode.
- At the end of the four-month maintenance mode period, the minor release will be considered EOL (end of life) and cherry-picks to the associated branch are to be closed soon afterward.

The table below documents support matrix for Cluster API provider vSphere versions (older versions omitted).

| Minor Release | Status                  | Supported Until (including maintenance mode)                                         |
|---------------|-------------------------|--------------------------------------------------------------------------------------|
| v1.16.x       | Standard support period | in maintenance mode when v1.18.0 will be released, EOL when v1.19.0 will be released |
| v1.15.x       | Standard support period | in maintenance mode when v1.17.0 will be released, EOL when v1.18.0 will be released |
| v1.14.x       | Maintenance mode        | Maintenance mode since v1.16.0 release date, EOL when v1.17.0 will be released       |
| v1.13.x       | EOL                     | EOL since v1.16.0 release date                                                       |
| v1.12.x       | EOL                     | EOL since v1.15.0 release date                                                       |
| v1.11.x       | EOL                     | EOL since v1.14.0 release date                                                       |

## Compatibility with Cluster API

Only a limited set of combinations between releases of CAPV and Cluster API are tested on CI and thus actively
supported:

| CAPV version | CAPI version used in CI |
|--------------|:-----------------------:|
| CAPV v1.17   |       CAPI v1.14        |
| CAPV v1.16   |       CAPI v1.13        |
| CAPV v1.15   |       CAPI v1.12        |
| CAPV v1.14   |       CAPI v1.11        |
| CAPV v1.13   |       CAPI v1.10        |

Other CAPV/CAPI version combinations might work, but it is up to the users to validate those combinations; 
also please be aware that in case of a CAPI/CAPV combination not listed above, support will be provided at best effort.

Following info should be used when identifying possible CAPV/CAPI version combinations outside the table above:

Compatibility between CAPV and Cluster API depends on which Cluster API contract versions is supported by CAPI and by
which contract version is implemented in CAPV.

The following table documents the Cluster API contract versions implemented by supported CAPV versions:

| CAPV version | Implements Cluster API contract version |
|--------------|:---------------------------------------:|
| CAPV v1.17   |                 v1beta2                 |
| CAPV v1.16   |                 v1beta2                 |
| CAPV v1.15   |                 v1beta1                 |
| CAPV v1.14   |                 v1beta1                 |
| CAPV v1.13   |                 v1beta1                 |

Contract versions supported by Cluster API are documented in the [versions](https://cluster-api.sigs.k8s.io/reference/versions#cluster-api-release-vs-contract-versions) page in the [Cluster API book](https://cluster-api.sigs.k8s.io).

Additionally, each version of CAPV requires a specific API version to served by Cluster API, and this
further reduces the number of valid combinations.

The table below documents CAPI API version required by supported CAPV versions.

| CAPV version | Requires CAPI API version |
|--------------|:-------------------------:|
| CAPV v1.17   |          v1beta2          |
| CAPV v1.16   |          v1beta2          |
| CAPV v1.15   |          v1beta2          |
| CAPV v1.14   |          v1beta2          |
| CAPV v1.13   |          v1beta1          |

API versions supported by Cluster API are documented in the [versions](https://cluster-api.sigs.k8s.io/reference/versions#cluster-api-release-vs-api-versions) page in the [Cluster API book](https://cluster-api.sigs.k8s.io).

## Compatibility with Kubernetes Versions

CAPV doesn't dictate supported K8s versions, but in practice compatibility with Kubernetes versions is limited by the
compatibility matrix of the Cluster API version you are using in combination with CAPV.

Please refer to Cluster API 's [Kubernetes versions matrix](https://cluster-api.sigs.k8s.io/reference/versions#kubernetes-versions-support).

## Kubernetes versions with published OVAs

**Note:** These OVAs are **not updated for security fixes** and it is recommended to always use the latest 
versions for distribution packages and patch version for the Kubernetes version you wish to run. For
production-like environments, it is highly recommended to build and use your own custom images.

**Note:** We recently moved the OVAs from the community GCP project to [Github releases](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases). Going forward new OVAs will only be uploaded to a dedicated Github release with the name `templates/<Kubernetes Version>`.

**Note:** Big OVAs will be split into multiple parts. To use them please download all parts and use `cat part1 part2 part3 > out.ova` to join them again.

| Kubernetes | Ubuntu 18.04 | Ubuntu 20.04 | Ubuntu 22.04 | Ubuntu 24.04 | Photon 3 | Photon 5 | Flatcar Stable |
|:-----------|:------------:|:------------:|:------------:|:------------:|:--------:|:--------:|:--------------:|
| [v1.36.0]  |              |              |      ✓       |      ✓       |          |          |       ✓        |
| [v1.35.0]  |              |              |      ✓       |      ✓       |          |    ✓     |       ✓        |
| [v1.34.0]  |              |              |      ✓       |      ✓       |          |    ✓     |       ✓        |
| [v1.33.0]  |              |              |      ✓       |      ✓       |          |    ✓     |       ✓        |
| [v1.32.0]  |              |              |      ✓       |      ✓       |          |    ✓     |       ✓        |
| [v1.31.0]  |              |              |      ✓       |      ✓       |          |    ✓     |       ✓        |
| [v1.30.0]  |              |              |      ✓       |              |          |    ✓     |       ✓        |
| [v1.29.0]  |              |              |      ✓       |              |    ✓     |    ✓     |       ✓        |
| [v1.28.0]  |      ✓       |      ✓       |      ✓       |              |    ✓     |          |       ✓        |
| [v1.27.3]  |      ✓       |      ✓       |      ✓       |              |    ✓     |          |       ✓        |
| [v1.26.2]  |      ✓       |      ✓       |              |              |    ✓     |          |       ✓        |
| [v1.25.7]  |      ✓       |      ✓       |              |              |    ✓     |          |       ✓        |
| [v1.24.11] |      ✓       |      ✓       |              |              |    ✓     |          |       ✓        |

[v1.36.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.36.0
[v1.35.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.35.0
[v1.34.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.34.0
[v1.33.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.33.0
[v1.32.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.32.0
[v1.31.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.31.0
[v1.30.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.30.0
[v1.29.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.29.0
[v1.28.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.28.0
[v1.27.3]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.27.3
[v1.26.2]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.26.2
[v1.25.7]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.25.7
[v1.24.11]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.24.11

A full list of the published machine images for CAPV can be found by [searching for releases](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases?q=templates%2F&expanded=true)
on the repository having the prefix `templates/` and taking a look at the available assets.

## Documentation

Further documentation is available in the `/docs` directory.

[vSphere Custom Resource Definitions][vsphere_custom_resource_definitions]

[Cluster API Custom Resource Definitions][capi_custom_resource_definitions]

  ## Getting involved and contributing

Are you interested in contributing to cluster-api-provider-vsphere? We, the maintainers and community, would love your suggestions, contributions, and help! Also, the maintainers can be contacted at
any time to learn more about how to get involved.

In the interest of getting more new people involved we tag issues with [`good first issue`][good_first_issue]. These are typically issues that have smaller scope but are good ways to start to get
acquainted with the codebase.

We also encourage ALL active community participants to act as if they are maintainers, even if you don't have "official" write permissions. This is a community effort, we are here to serve the
Kubernetes community. If you have an active interest and you want to get involved, you have real power! Don't assume that the only people who can get things done around here are the "maintainers".

We also would love to add more "official" maintainers, so show us what you can do!

This repository uses the Kubernetes bots. See a full list of the commands [here][prow].

## Code of conduct

Participating in the project is governed by the Kubernetes code of conduct. Please take some time to read the [code of conduct document][code_of_conduct].

### Implementer office hours

- Bi-weekly on [Zoom][zoom_meeting] on Wednesdays @ 02:30 PM Central European Time.
- Previous meetings: \[ [notes][meeting_notes] \]

### Other ways to communicate with the contributors

Please check in with us in the [#cluster-api-vsphere][slack] channel on Slack or email us at our [mailing list][mailing_list]

## Github issues

### Bugs

If you think you have found a bug please follow the instructions below.

- Please spend a small amount of time giving due diligence to the issue tracker. Your issue might be a duplicate.
- Get the logs from the cluster controllers. Please paste this into your issue.
- Follow the helpful tips provided in the [troubleshooting document][troubleshooting] as needed.
- Open a [new issue][new_issue].
- Remember that users might be searching for your issue in the future, so please give it a meaningful title to help others.
- Feel free to reach out to the cluster-api community on the [kubernetes slack][slack_info].

### Tracking new features

We also use the issue tracker to track features. If you have an idea for a feature, or think you can help CAPV become even more awesome follow the steps below.

- Open a [new issue][new_issue].
- Remember that users might be searching for your issue in the future, so please give it a meaningful title to help others.
- Clearly define the use case, using concrete examples. EG: I type `this` and cluster-api-provider-vsphere does `that`.
- Some of our larger features will require some design. If you would like to include a technical design for your feature please include it in the issue.
- After the new feature is well understood, and the design agreed upon, we can start coding the feature. We would love for you to code it. So please open up a **WIP** *(work in progress)* pull
  request, and happy coding.

<!-- References -->

[cluster_api]: https://github.com/kubernetes-sigs/cluster-api

[code_of_conduct]: https://git.k8s.io/community/code-of-conduct.md

[good_first_issue]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22

[kops]: https://github.com/kubernetes/kops

[kubicorn]: https://github.com/kubicorn/kubicorn

[mailint_list]: https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle

[meeting_notes]: https://docs.google.com/document/d/15CD2VOdkCAEcq2mm5FVoPO8M4-0a2SA2ajHLFBYqz7c/edit?usp=sharing

[new_issue]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/issues/new

[prow]: https://prow.k8s.io/command-help?repo=kubernetes-sigs%2Fcluster-api-provider-vsphere

[slack]: https://kubernetes.slack.com/messages/CKFGK3SSD

[slack_info]: https://github.com/kubernetes/community/tree/master/communication#communication

[troubleshooting]: ./docs/troubleshooting.md

[zoom_meeting]: https://zoom.us/j/92253194848?pwd=cVVVNDMxeTl1QVJPUlpvLzNSVU1JZz09

[vsphere_custom_resource_definitions]: https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api-provider-vsphere

[capi_custom_resource_definitions]: https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api

<!-- markdownlint-disable-file MD033 -->
