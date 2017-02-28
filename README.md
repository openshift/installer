# Tectonic Platform SDK

![Lifecycle Prototype](https://img.shields.io/badge/Lifecycle-Prototype-f4cccc.svg)

The Tectonic Platform SDK provides pre-built recipes to help users create the underlying compute infrastructure for a [Self-Hosted Kubernetes Cluster](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/self-hosted-kubernetes.md) ([vid](https://coreos.com/blog/self-hosted-kubernetes.html)) using [Hashicorp Terraform](https://terraform.io), [bootkube](https://github.com/kubernetes-incubator/bootkube), and supporting tooling.

The goal is to provide well-tested defaults that can be customized for various environments and plugged into other systems.

The unique power of Self-Hosted Kubernetes is that it cleanly separates out the infrastructure from Kubernetes enabling this separation of concerns:

![](http://i.imgur.com/Gd9W7RR.gif)

## Getting Started

TODO: @s-urbaniak will document how to hack this with the Tectonic Installer today: https://github.com/coreos-inc/tectonic-platform-sdk/issues/6#issuecomment-283003100

## Roadmap

This is an unprioritized list of future items the team would like to tackle:

- Run [Kubernetes e2e tests](https://github.com/coreos-inc/tectonic-platform-sdk/issues/6) over repo automatically
- Build a tool to walk the Terraform graph and warn if cluster won't comply with [Generic Platform](https://github.com/coreos-inc/tectonic-platform-sdk/blob/master/Documentation/generic-platform.md)
- Additional platforms like Azure, VMware, Google Cloud, etc
- Create a spec for generic and platform specific Terraform Variable files
- Document how to customize each of the platforms
- Create a tool to verify Terraform Variable files
- Deploy with other self-hosted tools like kubeadm
- Terraform plugin and integration with [matchbox](https://github.com/coreos/matchbox) for bare metal
