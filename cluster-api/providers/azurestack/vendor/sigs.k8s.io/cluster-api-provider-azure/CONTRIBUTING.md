# Contributing Guidelines
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Contributor License Agreements](#contributor-license-agreements)
- [Finding Things That Need Help](#finding-things-that-need-help)
- [Contributing a Patch](#contributing-a-patch)
- [Backporting a Patch](#backporting-a-patch)
  - [Merge Approval](#merge-approval)
  - [Google Doc Viewing Permissions](#google-doc-viewing-permissions)
  - [Issue and Pull Request Management](#issue-and-pull-request-management)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

Read the following guide if you're interested in contributing to cluster-api.

## Contributor License Agreements

We'd love to accept your patches! Before we can take them, we have to jump a couple of legal hurdles.

Please fill out either the individual or corporate Contributor License Agreement (CLA). More information about the CLA and instructions for signing it [can be found here](https://github.com/kubernetes/community/blob/master/CLA.md).

***NOTE***: Only original source code from you and other people that have signed the CLA can be accepted into the repository.

## Finding Things That Need Help

If you're new to the project and want to help, but don't know where to start, we have a semi-curated list of issues that should not need deep knowledge of the system. [Have a look and see if anything sounds interesting](https://github.com/kubernetes-sigs/cluster-api-provider-azure/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22). Alternatively, read some of the docs on other controllers and try to write your own, file and fix any/all issues that come up, including gaps in documentation!

## Contributing a Patch

1. If you haven't already done so, sign a Contributor License Agreement (see details above).
2. Fork the desired repo, develop and test your code changes.
    1. See the [Development Guide](https://capz.sigs.k8s.io/developers/development.html) for more instructions on setting up your environment and testing changes locally.
3. Submit a pull request.
    1. All PRs should be labeled with one of the following kinds
         - `/kind feature` for PRs related to adding new features/tests
         - `/kind bug` for PRs related to bug fixes and patches
         - `/kind api-change` for PRs related to adding, removing, or otherwise changing an API
         - `/kind cleanup` for PRs related to code refactoring and cleanup
         - `/kind deprecation` for PRs related to a feature/enhancement marked for deprecation.
         - `/kind design` for PRs related to design proposals
         - `/kind documentation` for PRs related to documentation
         - `/kind failing-test` for PRs related to a consistently or frequently failing test.
         - `/kind flake` for PRs related to a flaky test.
         - `/kind other` for PRs related to updating dependencies, minor changes or other
     2. If the PR requires additional action from users switching to a new release, include the string "action required" in the PR release-notes.
     3. All code changes must be covered by unit tests and E2E tests.
     4. All new features should come with user documentation.
 4. Once the PR has been reviewed and is ready to be merged, commits should be [squashed](https://github.com/kubernetes/community/blob/master/contributors/guide/github-workflow.md#squash-commits).
    1. Ensure that commit message(s) are be meaningful and commit history is readable.

All changes must be code reviewed. Coding conventions and standards are explained in the official [developer docs](https://github.com/kubernetes/community/tree/master/contributors/devel). Expect reviewers to request that you avoid common [go style mistakes](https://github.com/golang/go/wiki/CodeReviewComments) in your PRs.

In case you want to run our E2E tests locally, please refer to [Testing](https://capz.sigs.k8s.io/developers/development.html#submitting-prs-and-testing) guide. An overview of our e2e-test jobs (and also all our other jobs) can be found in [Jobs](https://capz.sigs.k8s.io/developers/jobs.html).

## Backporting a Patch

Cluster API ships older versions through `release-X.X` branches, usually backports are reserved to critical bug-fixes.
Some release branches might ship with both Go modules and dep (e.g. `release-0.1`), users backporting patches should always make sure
that the vendored Go modules dependencies match the Gopkg.lock and Gopkg.toml ones by running `dep ensure`

### Merge Approval

Cluster API maintainers may add "LGTM" (Looks Good To Me) or an equivalent comment to indicate that a PR is acceptable. Any change requires at least one LGTM.  No pull requests can be merged until at least one Cluster API maintainer signs off with an LGTM.

### Google Doc Viewing Permissions

To gain viewing permissions to google docs in this project, please join either the [kubernetes-dev](https://groups.google.com/forum/#!forum/kubernetes-dev) or [kubernetes-sig-cluster-lifecycle](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle) google group.

### Issue and Pull Request Management

Anyone may comment on issues and submit reviews for pull requests. However, in
order to be assigned an issue or pull request, you must be a member of the
[Kubernetes SIGs](https://github.com/kubernetes-sigs) GitHub organization.

If you are a Kubernetes GitHub organization member, you are eligible for
membership in the Kubernetes SIGs GitHub organization and can request
membership by [opening an issue](https://github.com/kubernetes/org/issues/new?template=membership.md&title=REQUEST%3A%20New%20membership%20for%20%3Cyour-GH-handle%3E)
against the kubernetes/org repo.

However, if you are a member of any of the related Kubernetes GitHub
organizations but not of the Kubernetes org, you will need explicit sponsorship
for your membership request. You can read more about Kubernetes membership and
sponsorship [here](https://github.com/kubernetes/community/blob/master/community-membership.md).

Cluster API maintainers can assign you an issue or pull request by leaving a
`/assign <your Github ID>` comment on the issue or pull request.

## Mentorship

- [Mentoring Initiatives](https://git.k8s.io/community/mentoring) - We have a diverse set of mentorship programs available that are always looking for volunteers!

## Contact Information

- [Slack channel](https://kubernetes.slack.com/messages/cluster-api-azure) - [Mailing list](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle)
