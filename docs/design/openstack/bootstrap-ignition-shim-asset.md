---
title: Bootstrap Ignition Shim Asset
authors:
  - "@iamemilio"
reviewers:
  - TBD
  - "@abhinavdahiya"
  - "@wking"
approvers:
  - TBD
  - "@abhinavdahiya"
  - "@wking"
creation-date: 2019-09-04
last-updated: 2019-09-04
status: implementable
---

# Bootstrap Ignition Shim Asset

## Release Signoff Checklist

- [ ] Enhancement is `implementable`
- [ ] Design details are appropriately documented from clear requirements
- [ ] Test plan is defined

## Summary

There are a few platforms that use an Ignition shim for their bootstrap machine.
Rather than generate this shim in Terraform, platforms should be able to **Opt-In**
to generate a bootstrap Ignition shim in a Go based installer asset. The shim serves
as a way to avoid hitting metadata size limits while still providing the machine with
the data it needs to succesfully fetch the main Ignition config.

## Motivation

The motivation for this is largely due to the sluggish pace of updates to the
Terraform Ignition Provider, which currently only supports ignition v2.1.0.
This has prevented the OpenStack team from fixing bugs related to CA Cert
Bundles not being trusted on the bootstrap machine. On top of this, it would
consolidate the generation of Ignition configs to one region of the installer,
reducing the overall complexity.

### Goals

Our goal is to create an installer asset, written in Go, that creates an Ignition
shim for the bootstrap node. This shim should contain the minimum amount of data
required to succesfully fetch the bootstrap Ignition config from its endpoint.
All further customizations should be added to the bootstrap Ignition config, not
the bootstrap shim. The design of this asset should be generic so that platforms can
implement customizations to it, like they do in the boostrap ignition asset.

### Risks and Mitigations

- This implementation has to be generic enough to be useable by the other platforms
that rely on an ignition shim. To mitigate this, developers from those teams will
need to be included in the design and review process of this feature in order to
make sure that their needs are met.

- Some values, such as the Swift temp url in OpenStack, are generated in Terraform,
and are therefore unable to be moved to the new installer asset.

## Design Details

### Test Plan

Like the other Ignition assets in the installer, this should be developed along side
a set of unit tests.

## Alternatives

- We could invest more resources in keeping the Terraform Ignition Provider up to date.
