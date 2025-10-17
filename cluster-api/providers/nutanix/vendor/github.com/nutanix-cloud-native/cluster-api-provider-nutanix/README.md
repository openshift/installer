# Kubernetes cluster-api infrastructure provider Nutanix Cloud Infrastucture

---

[![Go Report Card](https://goreportcard.com/badge/github.com/nutanix-cloud-native/cluster-api-provider-nutanix)](https://goreportcard.com/report/github.com/nutanix-cloud-native/cluster-api-provider-nutanix)
[![CI](https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix/actions/workflows/build-dev.yaml/badge.svg)](https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix/actions/workflows/build-dev.yaml)
[![Release](https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix/actions/workflows/release.yaml/badge.svg)](https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix/actions/workflows/release.yaml)

[![release](https://img.shields.io/github/release-pre/nutanix-cloud-native/cluster-api-provider-nutanix.svg)](https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix/blob/master/LICENSE)
![Proudly written in Golang](https://img.shields.io/badge/written%20in-Golang-92d1e7.svg)
[![Releases](https://img.shields.io/github/downloads/nutanix-cloud-native/cluster-api-provider-nutanix/total.svg)](https://github.com/nutanix-cloud-native/cluster-api-provider-nutanix/releases)

---
## What is the Cluster API Provider Nutanix Cloud Infrastucture
The [Cluster API](https://github.com/kubernetes-sigs/cluster-api) brings declarative, Kubernetes-style APIs to cluster creation, configuration and management. Cluster API Provider for Nutanix Cloud Infrastructure is a concrete implementation of Cluster API for Nutanix Cloud Infrastructure.

The API itself is shared across multiple cloud providers allowing for true Nutanix hybrid cloud deployments of Kubernetes. 

## How to Deploy a Kubernetes Cluster on Nutanix Cloud Infrastucture
Check out the [getting started guide](https://opendocs.nutanix.com/capx/latest/getting_started/) for launching a cluster on Nutanix Cloud Infrastructure.

## Compatibility with Prism Central & Prism Element

| CAPX Version | Min. Prism Central Version | Min. Prism Element Version |
|--------------|----------------------------|----------------------------|
| 1.5.x        | pc2024.1+                  | 6.5+                       |
| 1.6.x        | pc2024.2+                  | 6.5+                       |

## Documentation
Visit the `Cluster API Provider: Nutanix (CAPX)` section on [opendocs.nutanix.com](https://opendocs.nutanix.com/) for all documentation related to CAPX.

## Development
The [Development Workflow](./docs/developer_workflow.md) page explains how to build and test CAPX from source.

## Contributing
See the [contributing docs](CONTRIBUTING.md).

## Support
### Community Plus

This code is developed in the open with input from the community through issues and PRs. A Nutanix engineering team serves as the maintainer. Documentation is available in the project repository.

Issues and enhancement requests can be submitted in the [Issues tab of this repository](../../issues). Please search for and review the existing open issues before submitting a new issue.

## License
The project is released under version 2.0 of the [Apache license](http://www.apache.org/licenses/LICENSE-2.0).
