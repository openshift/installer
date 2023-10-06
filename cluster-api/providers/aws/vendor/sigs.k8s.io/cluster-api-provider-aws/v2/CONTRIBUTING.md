# Contributing guidelines

## Sign the CLA

Kubernetes projects require that you sign a Contributor License Agreement (CLA) before we can accept your pull requests.  Please see https://git.k8s.io/community/CLA.md for more info

### Contributing A Patch

1. Submit an issue describing your proposed change to the repo in question.
2. The [repo owners](OWNERS) will respond to your issue promptly.
3. If your proposed change is accepted, and you haven't already done so, sign a Contributor License Agreement (see details above).
4. Fork the desired repo, develop and test your code changes. 
> See the [developer guide](https://cluster-api-aws.sigs.k8s.io/development/development.html) on how to setup your development environment.
5. Submit a pull request.

### Contributer Ladder

We broadly follow the requirements from the [Kubernetes Community Membership](https://github.com/kubernetes/community/blob/master/community-membership.md).

> When making changes to **OWNER_ALIASES** please check that the **sig-cluster-lifecycle-leads**, **cluster-api-admins** and **cluster-api-maintainers** are correct.

#### Becoming a reviewer

If you would like to become a reviewer, then please ask one of the current maintainers.

We generally try to follow the [requirements for a reviewer](https://github.com/kubernetes/community/blob/master/community-membership.md#reviewer) from upstream Kubernetes. But if you feel that you don't full meet the requirements then reach out to us, they are not set in stone.

A reviewer can get PRs automatically assigned for review, and can `/lgtm` PRs.

To become a reviewer, ensure you are a member of the **kubernetes-sigs** Github organisation
following https://github.com/kubernetes/org/issues/new/choose.

The steps to add someone as a reviewer are:

- Add the GitHub alias to the **cluster-api-aws-reviewers** section of [OWNERS_ALIASES](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/OWNERS_ALIASES)
- Create a PR with the change that is held (i.e. by using `/hold`)
- Announce the change within the CAPA slack channel and as a PSA in the next CAPA office hours
- After 7 days of lazy consensus or after the next CAPA office hours (whichever is longer) the PR can be merged

#### Becoming a maintainer

If you have made significant contributions to Cluster API Provider AWS, a maintainer may nominate you to become a maintainer for the project.

We generally follow the [requirements for a approver](https://github.com/kubernetes/community/blob/master/community-membership.md#approver) from upstream Kubernetes. However, if you don't fully meet the requirements then a quorum of maintainers may still propose you if they feel you will make significant contributions.

Maintainers are able to approve PRs, as well as participate in release processes and have write access to the repo. **As a maintainer you will be expected to run the office hours, especially if no else wants to**.

Maintainers require membership of the **Kubernetes** Github organisation via
https://github.com/kubernetes/org/issues/new/choose

The steps to add someone as a reviewer are:

- Add the GitHub alias to the **cluster-api-aws-maintainers** and remove them from **cluster-api-aws-reviewers** sections of [OWNERS_ALIASES](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/OWNERS_ALIASES)
- Create a PR with the change that is held (i.e. by using `/hold`)
- Announce the change within the CAPA slack channel and as a PSA in the next CAPA office hours
- After 7 days of lazy consensus or after the next CAPA office hours (whichever is longer) the PR can be merged
- Open PR to add Github username to **cluster-api-provider-aws-maintainers**
to https://github.com/kubernetes/org/blob/main/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml
- Open PR to add Github username to https://github.com/kubernetes/test-infra/blob/master/config/jobs/kubernetes-sigs/cluster-api-provider-aws/OWNERS
- Open PR to add Github username to https://github.com/kubernetes/k8s.io/blob/main/k8s.gcr.io/images/k8s-staging-cluster-api-aws/OWNERS
- Open PR to add Google ID to the k8s-infra-staging-cluster-api-aws@kubernetes.io Google group in https://github.com/kubernetes/k8s.io/blob/main/groups/groups.yaml

#### Becoming a admin

After a period of time one of the existing CAPA or CAPI admins may propose you to become an admin of the CAPA project.

Admins have GitHub **admin** access to perform tasks on the repo.

The steps to add someone as an admin are:

- Add the GitHub alias to the **cluster-api-aws-admins** section of [OWNERS_ALIASES](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/OWNERS_ALIASES)
- Create a PR with the change that is held (i.e. by using `/hold`)
- Announce the change within the CAPA slack channel and as a PSA in the next CAPA office hours
- After 7 days of lazy consensus or after the next CAPA office hours (whichever is longer) the PR can be merged
- Open PR to add Github username to **cluster-api-provider-aws-admins**
to https://github.com/kubernetes/org/blob/main/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml
