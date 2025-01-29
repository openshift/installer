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
- Choice of Linux distribution between Ubuntu 18.04 and CentOS 7 using VM Templates based on [OVA images](#Kubernetes-versions-with-published-OVAs).
- Deploys Kubernetes control planes into provided clusters on vSphere.
- Doesn't use SSH for bootstrapping nodes.
- Installs only the minimal components to bootstrap a control plane and workers.

------

## Compatibility with Cluster API and Kubernetes Versions

This provider's versions are compatible with the following versions of Cluster API:

|                      | Cluster API v1beta1 (v1.6) | Cluster API v1beta1 (v1.7) | Cluster API v1beta1 (v1.8) | Cluster API v1beta1 (v1.9) |
|----------------------|:--------------------------:|:--------------------------:|:--------------------------:|:--------------------------:|
| CAPV v1beta1 (v1.9)  |             ✓              |             x              |             x              |             x              |
| CAPV v1beta1 (v1.10) |             x              |             ✓              |             x              |             x              |
| CAPV v1beta1 (v1.11) |             x              |             x              |             ✓              |             x              |
| CAPV v1beta1 (v1.12) |             x              |             x              |             x              |             ✓              |

As CAPV doesn't dictate supported K8s versions, and it supports whatever CAPI supported, about the provider's compatibility with K8s versions, please refer
to [CAPI Supported Kubernetes Versions](https://cluster-api.sigs.k8s.io/reference/versions.html).

Basically:

- 4 Kubernetes minor releases for the management cluster (N - N-3)
- 6 Kubernetes minor releases for the workload cluster (N - N-5)

**NOTES:**
* We aim to cut a CAPV minor release approximately one week after the corresponding CAPI minor release is out.
* We aim to cut a CAPV minor or patch release with support for a new Kubernetes minor version approximately 3 business days after releases for CAPI and CPI that support the new Kubernetes version are available.

## Kubernetes versions with published OVAs

**Note:** These OVAs are **not updated for security fixes** and it is recommended to always use the latest 
versions for distribution packages and patch version for the Kubernetes version you wish to run. For
production-like environments, it is highly recommended to build and use your own custom images.

**Note:** We recently moved the OVAs from the community GCP project to [Github releases](https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases). Going forward new OVAs will only be uploaded to a dedicated Github release with the name `templates/<Kubernetes Version>`.

**Note:** Big OVAs will be split into multiple parts. To use them please download all parts and use `cat part1 part2 part3 > out.ova` to join them again.

| Kubernetes | Ubuntu 18.04 | Ubuntu 20.04 | Ubuntu 22.04 | Ubuntu 24.04 | Photon 3 | Photon 5 | Flatcar Stable |
|:-----------|:------------:|:------------:|:------------:|:------------:|:--------:|:--------:|:--------------:|
| [v1.24.11] |      ✓       |      ✓       |              |              |    ✓     |          |       ✓        |
| [v1.25.7]  |      ✓       |      ✓       |              |              |    ✓     |          |       ✓        |
| [v1.26.2]  |      ✓       |      ✓       |              |              |    ✓     |          |       ✓        |
| [v1.27.3]  |      ✓       |      ✓       |      ✓       |              |    ✓     |          |       ✓        |
| [v1.28.0]  |      ✓       |      ✓       |      ✓       |              |    ✓     |          |       ✓        |
| [v1.29.0]  |              |              |      ✓       |              |    ✓     |    ✓     |       ✓        |
| [v1.30.0]  |              |              |      ✓       |              |          |    ✓     |       ✓        |
| [v1.31.0]  |              |              |      ✓       |      ✓       |          |    ✓     |       ✓        |
| [v1.32.0]  |              |              |      ✓       |      ✓       |          |    ✓     |       ✓        |

[v1.24.11]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.24.11
[v1.25.7]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.25.7
[v1.26.2]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.26.2
[v1.27.3]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.27.3
[v1.28.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.28.0
[v1.29.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.29.0
[v1.30.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.30.0
[v1.31.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.31.0
[v1.32.0]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/releases/tag/templates/v1.32.0

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

[kubicorn]: http://kubicorn.io/

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
