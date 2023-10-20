# Kubernetes Cluster API Provider Azure

<p align="center">
<img src="https://github.com/kubernetes/kubernetes/raw/master/logo/logo.png"  width="100">
</p>
<p align="center">
<!-- go doc / reference card -->
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api-provider-azure">
<img src="https://godoc.org/sigs.k8s.io/cluster-api-provider-azure?status.svg"></a>
<!-- goreportcard badge -->
<a href="https://goreportcard.com/report/kubernetes-sigs/cluster-api-provider-azure">
<img src="https://goreportcard.com/badge/kubernetes-sigs/cluster-api-provider-azure"></a>
<!-- join kubernetes slack channel for cluster-api-azure-provider -->
<a href="https://kubernetes.slack.com/messages/CEX9HENG7">
<img src="https://img.shields.io/badge/join%20slack-%23cluster--api--azure-brightgreen"></a>
</p>

------

Kubernetes-native declarative infrastructure for Azure.

## What is the Cluster API Provider Azure

The [Cluster API][cluster_api] brings declarative, Kubernetes-style APIs to cluster creation, configuration and management.

The API itself is shared across multiple cloud providers allowing for true Azure
hybrid deployments of Kubernetes.

## Quick Start

Check out the [Cluster API Quick Start][quickstart] to create your first Kubernetes cluster on Azure using Cluster API.

## Flavors

See the [flavors documentation][flavors_doc] to know which cluster templates are provided by CAPZ.

## Getting Help

If you need help with CAPZ, please visit the [#cluster-api-azure][slack] channel on Slack, open a [GitHub issue](#github-issues), or join us at [Office Hours](#office-hours).

------

## Compatibility

### Cluster API Versions

Currently, CAPZ is compatible only with the `v1beta1` version of CAPI (v1.0.x). Support for `v1alpha3` (v0.3.x) and `v1alpha4` (v0.4.x) is deprecated and has been removed.

### Kubernetes Versions

The Azure provider is able to install and manage the [versions of Kubernetes supported by the Cluster API (CAPI) project](https://cluster-api.sigs.k8s.io/reference/versions.html#supported-kubernetes-versions).

#### Managed Clusters (AKS)

Managed Clusters (AKS) follow their own [Kubernetes version support policy](https://learn.microsoft.com/azure/aks/supported-kubernetes-versions?tabs=azure-cli#kubernetes-version-support-policy). Please use the Azure portal or CLI to [find the versions supported in your cluster's location](https://learn.microsoft.com/azure/aks/supported-kubernetes-versions?tabs=azure-cli#azure-portal-and-cli-versions).

For more information on Kubernetes version support, see the [Cluster API book](https://cluster-api.sigs.k8s.io/reference/versions.html).

------

## Documentation

Please see our [Book](https://capz.sigs.k8s.io) for in-depth user documentation.

Additional docs can be found in the `/docs` directory, and the [index is here](https://github.com/kubernetes-sigs/cluster-api-provider-azure/blob/main/docs/README.md).

## Getting involved and contributing

Are you interested in contributing to cluster-api-provider-azure? We, the
maintainers and community, would love your suggestions, contributions, and help!
Also, the maintainers can be contacted at any time to learn more about how to get
involved.

To set up your environment checkout the [development guide](https://capz.sigs.k8s.io/developers/development.html).

In the interest of getting more new people involved, we tag issues with
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

This repository uses the Kubernetes bots.  See a full list of the commands [here][prow].

### Office hours

The community holds office hours every week, with sessions open to all users and
developers.

Office hours are hosted on a zoom video chat every Thursday
at 09:00 (PT) / 12:00 (ET) [Convert to your timezone](https://www.thetimezoneconverter.com/?t=09:00&tz=PT%20%28Pacific%20Time%29)
and are published on the [Kubernetes community meetings calendar][gcal]. Please add your questions or ideas to [the agenda][capz_agenda].

### Other ways to communicate with the contributors

Please check in with us in the [#cluster-api-azure][slack] channel on Slack.

## Github issues

### Bugs

If you think you have found a bug please follow the instructions below.

- Please spend a small amount of time giving due diligence to the issue tracker. Your issue might be a duplicate.
- Get the logs from the cluster controllers. Please paste this into your issue.
- Open a [bug report][bug_report].
- Remember users might be searching for your issue in the future, so please give it a meaningful title to help others.
- Feel free to reach out to the cluster-api community on [kubernetes slack][slack_info].

### Tracking new features

We also use the issue tracker to track features. If you have an idea for a feature, or think you can help Cluster API Provider Azure become even more awesome, then follow the steps below.

- Open a [feature request][feature_request].
- Remember users might be searching for your issue in the future, so please
  give it a meaningful title to help others.
- Clearly define the use case, using concrete examples. EG: I type `this` and
  cluster-api-provider-azure does `that`.
- Some of our larger features will require some design. If you would like to
  include a technical design for your feature please include it in the issue.
- After the new feature is well understood, and the design agreed upon we can
  start coding the feature. We would love for you to code it. So please open
  up a **WIP** *(work in progress)* pull request, and happy coding.

<!-- References -->

[slack]: https://kubernetes.slack.com/messages/CEX9HENG7
[good_first_issue]: https://github.com/kubernetes-sigs/cluster-api-provider-azure/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22
[gcal]: https://calendar.google.com/calendar/embed?src=cgnt364vd8s86hr2phapfjc6uk%40group.calendar.google.com
[prow]: https://go.k8s.io/bot-commands
[bug_report]: https://github.com/kubernetes-sigs/cluster-api-provider-azure/issues/new?template=bug_report.md
[feature_request]: https://github.com/kubernetes-sigs/cluster-api-provider-azure/issues/new?template=feature_request.md
[slack_info]: https://github.com/kubernetes/community/tree/master/communication#slack
[cluster_api]: https://github.com/kubernetes-sigs/cluster-api
[quickstart]: https://cluster-api.sigs.k8s.io/user/quick-start.html
[flavors_doc]: https://github.com/kubernetes-sigs/cluster-api-provider-azure/blob/main/templates/flavors/README.md
[capz_agenda]: https://bit.ly/k8s-capz-agenda
