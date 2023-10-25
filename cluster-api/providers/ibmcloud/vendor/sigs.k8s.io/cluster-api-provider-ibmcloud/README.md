[![LICENSE](https://img.shields.io/badge/license-apache2.0-green.svg)](https://sigs.k8s.io/cluster-api-provider-ibmcloud/blob/master/LICENSE)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/kubernetes-sigs/cluster-api-provider-ibmcloud?label=version)
[![Go Report Card](https://goreportcard.com/badge/sigs.k8s.io/cluster-api-provider-ibmcloud)](https://goreportcard.com/report/sigs.k8s.io/cluster-api-provider-ibmcloud)

# Kubernetes Cluster API Provider IBM Cloud

<a href="https://github.com/kubernetes-sigs/cluster-api"><img src="https://github.com/kubernetes/kubernetes/raw/master/logo/logo.png"  width="100"></a><a href="https://www.ibm.com/cloud/"><img hspace="90px" src="./docs/images/ibm-cloud.svg" alt="Powered by IBM Cloud" height="100"></a>

------

This repository hosts a concrete implementation of an IBM Cloud provider for the [cluster-api project](https://github.com/kubernetes-sigs/cluster-api).

## What is the Cluster API Provider IBM Cloud

The [Cluster API](https://github.com/kubernetes-sigs/cluster-api) brings declarative, Kubernetes-style APIs to cluster creation, configuration and management. The API itself is shared across multiple cloud providers allowing for true IBM Cloud hybrid deployments of Kubernetes.

## Documentation

Please see our [book](https://cluster-api-ibmcloud.sigs.k8s.io) for in-depth documentation.

## Launching a Kubernetes cluster on IBM Cloud

Check out the [Cluster API IBM Cloud Quick Start](https://cluster-api-ibmcloud.sigs.k8s.io/getting-started.html) for launching a
cluster on IBM Cloud.

## Compatibility with Cluster API and Kubernetes Versions

This provider's versions are compatible with the following versions of Cluster API:

|                                            | v1alpha4 (v0.4) | v1beta1 (v1.x) |
|--------------------------------------------|-----------------|----------------|
| IBMCloud Provider v1alpha4 (v0.1.x)        | ✓               |                |
| IBMCloud Provider v1beta1 (v0.2.x, v0.3.x) |                 | ✓              |
| IBMCloud Provider v1beta2 (v0.4.x, v0.5.x, main)   |                 | ✓              |


(See [Kubernetes support matrix][cluster-api-supported-v] of Cluster API versions).

<!-- ANCHOR: Community -->

## Community, discussion, contribution, and support

If you have questions or want to get the latest project news, you can connect with us in the following ways:

- Chat with us on the Kubernetes [Slack](http://slack.k8s.io/) in the [#cluster-api-ibmcloud][slack] channel
- Subscribe to the [SIG Cluster Lifecycle](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle) Google Group for access to documents and calendars
- Join our Weekly meeting sessions where we share the latest project news, demos, answer questions, and triage issues
    - Biweekly on Friday @ 09:00 IST on [Zoom][zoomMeeting]
    - Previous meetings: \[ [notes][notes] \]

Pull Requests and feedback on issues are very welcome!
See the [issue tracker] if you're unsure where to start, especially the [Good first issue] and [Help wanted] tags, and
also feel free to reach out to discuss.

See also our [contributor guide](CONTRIBUTING.md) and the Kubernetes [community page] for more details on how to get involved.

[slack]: https://kubernetes.slack.com/messages/C02F4CX3ALF
[zoomMeeting]: https://zoom.us/j/9392903494
[notes]: https://docs.google.com/document/d/1oWnqXy1VFv0E3kovQoZfS6IlVP0L4eaQsN-2HYC_6_A/edit
[issue tracker]: https://github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/issues
[Good first issue]: https://github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22
[Help wanted]: https://github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/issues?utf8=%E2%9C%93&q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22+
[community page]: https://kubernetes.io/community
[cluster-api-supported-v]: https://cluster-api.sigs.k8s.io/reference/versions.html

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

<!-- ANCHOR_END: Community -->
