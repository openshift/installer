# Kubernetes Cluster API Provider OpenStack

[![Go Report Card](https://goreportcard.com/badge/kubernetes-sigs/cluster-api-provider-openstack)](https://goreportcard.com/report/kubernetes-sigs/cluster-api-provider-openstack)

------

Kubernetes-native declarative infrastructure for OpenStack.

For documentation, see the [Cluster API Provider OpenStack book](https://cluster-api-openstack.sigs.k8s.io/).

## What is the Cluster API Provider OpenStack

The [Cluster API][cluster_api] brings
declarative, Kubernetes-style APIs to cluster creation, configuration and
management.

The API itself is shared across multiple cloud providers allowing for true OpenStack
hybrid deployments of Kubernetes. It is built atop the lessons learned from
previous cluster managers such as [kops][kops] and
[kubicorn][kubicorn].


## Launching a Kubernetes cluster on OpenStack

- Check out the [Cluster API Quick Start](https://cluster-api.sigs.k8s.io/user/quick-start.html) to create your first Kubernetes cluster on OpenStack using Cluster API. If you wish to use the external cloud provider, check out the [External Cloud Provider](docs/book/src/topics/external-cloud-provider.md) as well.

## Features

- Native Kubernetes manifests and API
- Choice of Linux distribution (as long as a current cloud-init is available)
- Support for single and multi-node control plane clusters
- Deploy clusters with and without LBaaS available (only cluster with LBaaS can be upgraded)
- Support for security groups
- cloud-init based nodes bootstrapping

------

## Compatibility with Cluster API and Kubernetes Versions

This provider's versions are compatible with the following versions of Cluster API:

|                                    | v1beta1 (v1.x) |
|------------------------------------| -------------- |
| OpenStack Provider v1alpha5 (v0.6) | ✓              |
| OpenStack Provider v1alpha6 (v0.7) | ✓              |
| OpenStack Provider v1alpha7 (v0.9) | ✓              |
| OpenStack Provider v1beta1         | ✓              |


This provider's versions are able to install and manage the following versions of Kubernetes:

|                                    | v1.25 | v1.26 | v1.27 | v1.28 |
|------------------------------------| ----- | ----- | ----- | ----- |
| OpenStack Provider v1alpha5 (v0.6) | ✓     | +     | +     | +     |
| OpenStack Provider v1alpha6 (v0.7) | ✓     | ✓     | ✓     | +     |
| OpenStack Provider v1alpha7 (v0.9) | +     | ✓     | ✓     | ★     |
| OpenStack Provider v1beta1         | +     | ✓     | ✓     | ★     |

This provider's versions are able to install Kubernetes to the following versions of OpenStack:

|                                    | Queens | Rocky | Stein | Train | Ussuri | Victoria | Wallaby | Xena | Yoga | Bobcat |
|------------------------------------| ------ | ----- | ----- | ----- | ------ | -------- | ------- | ---- | ---- | ------ |
| OpenStack Provider v1alpha5 (v0.6) | +      | +     | +     | +     | +      | ✓        | ✓       | ✓    | ✓    | ★      |
| OpenStack Provider v1alpha6 (v0.7) | +      | +     | +     | +     | +      | ✓        | ✓       | ✓    | ✓    | ★      |
| OpenStack Provider v1alpha7 (v0.9) |        | +     | +     | +     | +      | ✓        | ✓       | ✓    | ✓    | ★      |
| OpenStack Provider v1beta1         |        | +     | +     | +     | +      | ✓        | ✓       | ✓    | ✓    | ★      |

Test status:

- `★` currently testing
- `✓` previously tested
- `+` should work, but we weren't able to test it

Older versions may also work but we have not verified.

Each version of Cluster API for OpenStack will attempt to support two Kubernetes versions.

**NOTE:** As the versioning for this project is tied to the versioning of Cluster API, future modifications to this
policy may be made to more closely aligned with other providers in the Cluster API ecosystem.

**NOTE:** The minimum microversion of CAPI using nova is `2.60` now due to `server tags` support as well permitting `multiattach` volume types, see [code](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/c052e7e600f0e5ebddc839c08746bb636e79be87/pkg/cloud/services/compute/service.go#L38) for additional information.

**NOTE:** We require Keystone v3 for authentication.

------

## Development versions

ClusterAPI provider OpenStack images and manifests are published after every PR merge and once every day:

* With a Google Cloud account you can get a quick overview [here](https://console.cloud.google.com/storage/browser/artifacts.k8s-staging-capi-openstack.appspot.com/components)
* The manifests are available under:
  * [master/infrastructure-components.yaml](https://storage.googleapis.com/artifacts.k8s-staging-capi-openstack.appspot.com/components/master/infrastructure-components.yaml):
    latest build from the main branch, overwritten after every merge
  * e.g. [nightly_master_20210407/infrastructure-components.yaml](https://storage.googleapis.com/artifacts.k8s-staging-capi-openstack.appspot.com/components/nightly_master_20210407/infrastructure-components.yaml): build of the main branch from 7th April

These artifacts are published via Prow and Google Cloud Build. The corresponding job definitions can
be found [here](https://github.com/kubernetes/test-infra/blob/4d146721aaec27a3c93299956f8d64af2357d64a/config/jobs/image-pushing/k8s-staging-cluster-api.yaml).

------

## Operating system images

Note: Cluster API Provider OpenStack relies on a few prerequisites which have to be already
installed in the used operating system images, e.g. a container runtime, kubelet, kubeadm,.. .
Reference images can be found in [kubernetes-sigs/image-builder](https://github.com/kubernetes-sigs/image-builder/tree/master/images/capi). If it isn't possible to pre-install those
 prerequisites in the image, you can always deploy and execute some custom scripts
 through the [KubeadmConfig](https://github.com/kubernetes-sigs/cluster-api-bootstrap-provider-kubeadm).

------

## Documentation

Please see our [book](https://cluster-api-openstack.sigs.k8s.io) for in-depth documentation.

## Getting involved and contributing

Are you interested in contributing to cluster-api-provider-openstack? We, the
maintainers and community, would love your suggestions, contributions, and help!
Also, the maintainers can be contacted at any time to learn more about how to get
involved:

- via the [cluster-api-openstack channel on Kubernetes Slack][slack]
- via the [SIG-Cluster-Lifecycle Mailing List](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle).
- during our Office Hours
  - bi-weekly on Wednesdays @ 14:00 UTC on Zoom (link in meeting notes)
  - Previous meetings: [ [notes][office-hours-notes] | [recordings][office-hours-recordings] ]

In the interest of getting more new people involved we try to tag issues with
[`good first issue`][good_first_issue].
These are typically issues that have smaller scope but are good ways to start
to get acquainted with the codebase.

We also encourage ALL active community participants to act as if they are
maintainers, even if you don't have "official" write permissions. This is a
community effort, we are here to serve the Kubernetes community. If you have an
active interest and you want to get involved, you have real power! Don't assume
that the only people who can get things done around here are the "maintainers".

We also would love to add more "official" maintainers, so show us what you can
do!

This repository uses the Kubernetes bots. See a full list of the commands [here][prow].
Please also refer to the [Contribution Guide](CONTRIBUTING.md) and the [Development Guide](docs/book/src/development/development.md) for this project.

## Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

## Github issues

### Bugs

If you think you have found a bug please follow the instructions below.

- Please spend a small amount of time giving due diligence to the issue tracker. Your issue might be a duplicate.
- Get the logs from the cluster controllers. Please paste this into your issue.
- Open a [new issue][new_bug_issue].
- Remember that users might be searching for your issue in the future, so please give it a meaningful title to help others.
- Feel free to reach out to the Cluster API community on the [Kubernetes Slack][slack].

### Tracking new features

We also use the issue tracker to track features. If you have an idea for a feature, or think you can help Cluster API Provider OpenStack become even more awesome follow the steps below.

- Open a [new issue][new_feature_issue].
- Remember that users might be searching for your issue in the future, so please
  give it a meaningful title to help others.
- Clearly define the use case, using concrete examples.
- Some of our larger features will require some design. If you would like to
  include a technical design for your feature, please include it in the issue.
- After the new feature is well understood, and the design agreed upon, we can
  start coding the feature. We would love for you to code it. So please open
  up a **WIP** *(work in progress)* pull request, and happy coding.


<!-- References -->

[cluster_api]: https://github.com/kubernetes-sigs/cluster-api
[kops]: https://github.com/kubernetes/kops
[kubicorn]: http://kubicorn.io/
[slack]: https://kubernetes.slack.com/messages/cluster-api-openstack
[office-hours-notes]: https://docs.google.com/document/d/1hzi6nr04mhQYBKrwL2NDTNPvgI4RgO9a-gqmk31kXMA/edit
[office-hours-recordings]: https://www.youtube.com/playlist?list=PL69nYSiGNLP29D0nYgAGWt1ZFqS9Z7lw4
[good_first_issue]: https://github.com/kubernetes-sigs/cluster-api-provider-openstack/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22
[prow]: https://go.k8s.io/bot-commands
[new_bug_issue]: https://github.com/kubernetes-sigs/cluster-api-provider-openstack/issues/new?assignees=&labels=&template=bug_report.md
[new_feature_issue]: https://github.com/kubernetes-sigs/cluster-api-provider-openstack/issues/new?assignees=&labels=&template=feature_request.md
