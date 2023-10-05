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

### Becoming a reviewer

If you would like to become a reviewer, then please ask one of the maintainers.
There's no hard and defined limit as to who can become a reviewer, but a good
heuristic is 5 or more contributions. A reviewer can get PRs automatically assigned
for review, and can `/lgtm` PRs.

To become a reviewer, ensure you are a member of the kubernetes-sigs Github organisation
following https://github.com/kubernetes/org/issues/new/choose .

### Steps needed to become a maintainer
If you have made significant contributions to Cluster API
Provider AWS, a maintainer may nominate you to become a
maintainer, first by opening a PR to add you to the OWNERS_ALIASES file of the repository.

Maintainers are able to approve PRs, as well as participate
in release processes.

Maintainers require membership of the Kubernetes Github organisation via
https://github.com/kubernetes/org/issues/new/choose

The complete list of tasks required to set up maintainer status
follow:

* Open PR to add Github username to the OWNERS_ALIASES file under cluster-api-aws-maintainers
* Open PR to add Github username to cluster-api-provider-aws-admins and cluster-api-provider-aws-maintainers
to https://github.com/kubernetes/org/blob/main/config/kubernetes-sigs/sig-cluster-lifecycle/teams.yaml
* Open PR to add Github username to https://github.com/kubernetes/test-infra/blob/master/config/jobs/kubernetes-sigs/cluster-api-provider-aws/OWNERS
* Open PR to add Github username to https://github.com/kubernetes/k8s.io/blob/main/k8s.gcr.io/images/k8s-staging-cluster-api-aws/OWNERS
* Open PR to add Google ID to the k8s-infra-staging-cluster-api-aws@kubernetes.io Google group in https://github.com/kubernetes/k8s.io/blob/main/groups/groups.yaml
