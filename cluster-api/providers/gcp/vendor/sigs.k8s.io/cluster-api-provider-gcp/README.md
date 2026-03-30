<p align="center"><img alt="capi" src="https://github.com/kubernetes-sigs/cluster-api/raw/main/docs/book/src/images/introduction.svg" width="160x" /><img alt="capi" src="https://cloud.google.com/_static/cloud/images/favicons/onecloud/super_cloud.png" width="192x" /></p>
<p align="center"><a href="https://prow.k8s.io/?job=ci-cluster-api-provider-gcp-build">
<!-- prow build badge, godoc, and go report card-->
<img alt="Build Status" src="https://prow.k8s.io/badge.svg?jobs=ci-cluster-api-provider-gcp">
</a> <a href="https://godoc.org/sigs.k8s.io/cluster-api-provider-gcp"><img src="https://godoc.org/sigs.k8s.io/cluster-api-provider-gcp?status.svg"></a> <a href="https://goreportcard.com/report/sigs.k8s.io/cluster-api-provider-gcp"><img alt="Go Report Card" src="https://goreportcard.com/badge/sigs.k8s.io/cluster-api-provider-gcp" /></a></p>

----

# Kubernetes Cluster API Provider GCP

Kubernetes-native declarative infrastructure for GCP.

## What is the Cluster API Provider GCP?

The [Cluster API](https://github.com/kubernetes-sigs/cluster-api) brings declarative Kubernetes-style APIs to cluster creation, configuration and management. The API itself is shared across multiple cloud providers allowing for true Google Cloud hybrid deployments of Kubernetes.

## Documentation

Please see our [book](https://cluster-api-gcp.sigs.k8s.io/) for in-depth documentation.

## Quick Start

Checkout our [Cluster API Quick Start] to create your first Kubernetes cluster
on Google Cloud Platform using Cluster API.

----

## Support Policy

This provider's versions are compatible with the following versions of Cluster API:

|  | Cluster API `v1alpha3` (`v0.3.x`) | Cluster API `v1alpha4` (`v0.4.x`) | Cluster API `v1beta1` (`v1.0.x`) |
|---|---|---|---|
|Google Cloud Provider `v0.3.x` | ✓ |  |  |
|Google Cloud Provider `v0.4.x` |  | ✓ |  |
|Google Cloud Provider `v1.0.x` |  |  | ✓ |

This provider's versions are able to install and manage the following versions of Kubernetes:

|  | Google Cloud Provider `v0.3.x` | Google Cloud Provider `v0.4.x` | Google Cloud Provider `v1.0.x` |
|---|:---:|:---:|:---:|
| Kubernetes 1.15 |  |  |  |
| Kubernetes 1.16 | ✓ |  |  |
| Kubernetes 1.17 | ✓ | ✓ |  |
| Kubernetes 1.18 | ✓ | ✓ | ✓ |
| Kubernetes 1.19 | ✓ | ✓ | ✓ |
| Kubernetes 1.20 | ✓ | ✓ | ✓ |
| Kubernetes 1.21 |  | ✓ | ✓ |
| Kubernetes 1.22 |  |  | ✓ |
 
Each version of Cluster API for Google Cloud will attempt to support at least two versions of Kubernetes e.g., Cluster API for GCP `v0.1` may support Kubernetes 1.13 and Kubernetes 1.14.

**NOTE:** As the versioning for this project is tied to the versioning of Cluster API, future modifications to this policy may be made to more closely align with other providers in the Cluster API ecosystem. 

----

## Getting Involved and Contributing

Are you interested in contributing to cluster-api-provider-gcp? We, the maintainers 
and the community would love your suggestions, support and contributions! The maintainers
of the project can be contacted anytime to learn about how to get involved.

Before starting with the contribution, please go through [prerequisites] of the project.

To set up the development environment, checkout the [development guide].

In the interest of getting new people involved, we have issues marked as [`good first issue`][good_first_issue]. Although
these issues have a smaller scope but are very helpful in getting acquainted with the codebase.
For more, see the [issue tracker]. If you're unsure where to start, feel free to reach out to discuss.

See also: Our own [contributor guide] and the Kubernetes [community page].

We also encourage ALL active community participants to act as if they are maintainers, even if you don't have
'official' written permissions. This is a community effort and we are here to serve the Kubernetes community.
If you have an active interest and you want to get involved, you have real power!


### Office hours

- Join the [SIG Cluster Lifecycle](https://groups.google.com/a/kubernetes.io/g/sig-cluster-lifecycle) Google Group for access to documents and calendars.
- Participate in the conversations on [Kubernetes Discuss][kubernetes discuss]
- Provider implementers office hours (CAPI)
    - Weekly on Wednesdays @ 10:00 am PT (Pacific Time) on [Zoom](https://zoom.us/j/861487554)
    - Previous meetings: \[ [notes][notes] | [recordings][recordings] \]
- Cluster API Provider GCP office hours (CAPG)
    - Monthly on first Thursday @ 05:00 am PT (Pacific Time) on [Zoom](https://zoom.us/j/96963829102?pwd=WjBZcmwvZFZsUU93aVZieUk1L3FnZz09)
    - Previous meetings: [ [notes](http://bit.ly/k8s-capg-agenda)|[recordings](https://www.youtube.com/playlist?list=PL69nYSiGNLP29D0nYgAGWt1ZFqS9Z7lw4) ]

### Other ways to communicate with the contributors

Please check in with us in the [#cluster-api-gcp] on Slack. 

## Github Issues

### Bugs

If you think you have found a bug, please follow the instruction below.

- Please give a small amount of time giving due diligence to the issue tracker. Your issue might be a duplicate.
- Get the logs from the custom controllers and please paste them in the issue.
- Open a [bug report].
- Remember users might be searching for the issue in the future, so please make sure to give it a meaningful title to help others.
- Feel free to reach out to the community on slack.

### Tracking new feature

We also have an issue tracker to track features. If you think you have a feature idea, that could make Cluster API provider GCP become even more awesome, then follow these steps.

- Open a [feature request].
- Remember users might be searching for the issue in the future, so please make sure to give it a meaningful title to help others.
- Clearly define the use case with concrete examples. Example: type `this` and cluster-api-provider-gcp does `that`.
- Some of our larger features will require some design. If you would like to include a technical design in your feature, please go ahead.
- After the new feature is well understood and the design is agreed upon, we can start coding the feature. We would love for you to code it. So please open up a **WIP** *(work in progress)* PR and happy coding!

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct].

[Cluster API Quick Start]: https://cluster-api.sigs.k8s.io/user/quick-start.html
[prerequisites]: https://github.com/kubernetes-sigs/cluster-api-provider-gcp/blob/main/docs/book/src/topics/prerequisites.md
[development guide]: https://github.com/kubernetes-sigs/cluster-api-provider-gcp/blob/main/docs/book/src/developers/development.md
[good_first_issue]: https://github.com/kubernetes-sigs/cluster-api-provider-gcp/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22
[issue tracker]: https://github.com/kubernetes-sigs/cluster-api-provider-gcp/issues
[contributor guide]: CONTRIBUTING.md 
[community page]: https://kubernetes.io/community
[Kubernetes Code of Conduct]: code-of-conduct.md
[notes]: https://docs.google.com/document/d/1LdooNTbb9PZMFWy3_F-XAsl7Og5F2lvG3tCgQvoB5e4
[recordings]: https://www.youtube.com/playlist?list=PL69nYSiGNLP29D0nYgAGWt1ZFqS9Z7lw4
[#cluster-api-gcp]: https://sigs.k8s.io/cluster-api-provider-gcp
[bug report]: https://github.com/kubernetes-sigs/cluster-api-provider-gcp/issues/new?assignees=&labels=&template=bug_report.md
[feature request]: https://github.com/kubernetes-sigs/cluster-api-provider-gcp/issues/new?assignees=&labels=&template=feature_request.md
[kubernetes discuss]: https://groups.google.com/a/kubernetes.io/g/sig-cluster-lifecycle
