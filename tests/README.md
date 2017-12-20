# Tectonic Installer Tests


## Running basic tests on PRs

Our basic set of tests includes:
- Code linting
- UI tests
- Backend unit tests

They are run on **every** PR by default. Successful basic tests are required in
order to merge any PRs. Before starting to test your proposed changes, they are
temporarily merged into the target branch of the pull request.

### Actions required
- **none**

### FAQ
- *I am not part of the `coreos` or `coreos-inc` GitHub organization. Why are
  the tests not being executed on my changes?*

  Pull requests by external contributors need to be checked before they are
  tested by our Jenkins setup. Please ask a [maintainer](../MAINTAINERS) to mark
  your PR via commenting `ok to test`.

- *How do I retrigger the tests?*

  Comment with `ok to test` on the PR.


## Running GUI tests on PRs

The GUI tests include integration tests for the AWS and the Baremetal GUI
installer.

### Actions required
- Add the `run-gui-tests` GitHub label

### FAQ
- *I am not able to add labels, what should I do?*

  Please ask one of the repository [maintainers](../MAINTAINERS) to add the
  labels.

- *How do I retrigger the tests?*

  Comment with `ok to test` on the PR.

- *I forgot to add the GitHub labels. Can I add them after creating the PR?*

  Yes, just add the GitHub labels and comment `ok to test` on the PR.


## Running smoke / k8s-conformance tests on PRs

In addition to our basic set of tests we have smoke tests and the k8s upstream
conformance tests. These test the Tectonic installer on our supported platforms:
- AWS
- Azure
- Bare metal
- GCP

### Actions required
- Add the `run-smoke-tests` or/and the `run-conformance-tests` GitHub label
- Add the `platform/<xxx>` GitHub label for **each** platform you want to run
  tests against

### FAQ
- *I am not able to add labels, what should I do?*

  Please ask one of the repository [maintainers](../MAINTAINERS) to add the
  labels.

- *How do I retrigger the tests?*

  comment with `ok to test` on the PR.

- *I forgot to add the GitHub labels. Can I add them after creating the PR?*

  Yes, just add the GitHub labels and comment `ok to test` on the PR.

- *What if I only add the `run-smoke-tests`/`run-conformance-tests` GitHub
  label, but no `platform/<xxx>` label?*

  No smoke / conformance tests will be executed.

- *What if I trigger the tests twice in a small time frame?*

  Triggering the tests twice in a small time frame results in two test
  executions. The result of the most recent execution will be reported as a PR
  status in GitHub.

- *What can I do in case I run into test flakes continually?*

  1. Make sure the test failure is in **no** way related to your PR changes.
     Test your changes locally thoroughly.
  2. Document the flake in Jira in the `INST` project as *issue type* "bug" with the
     `flake` label.
  3. Get the approval of another person.
  4. Merge the PR.

## Nightly Runs

We use Jenkins to run our tests, and we have two Jenkins Jobs to run our `nightly` builds against master branch.
One job run all smoke tests for all platforms we have tests for.
The other job run the conformance tests for one or all platforms we support. Today this job is running
the conformance tests for `Bare Metal`

Those builds report the status to our internal chat tool.

## Running smoke tests / k8s-conformance tests locally

### 1. Expose environment variables

To run the smoke tests / conformance tests locally you need to set the following
environment variables:
``` bash
CLUSTER
TF_VAR_tectonic_license_path
TF_VAR_tectonic_pull_secret_path
TF_VAR_tectonic_base_domain
TF_VAR_tectonic_admin_email
TF_VAR_tectonic_admin_password
TECTONIC_TESTS_DONT_CLEAN_UP // If you want to keep the cluster alive after the tests
RUN_SMOKE_TESTS=true
RUN_CONFORMANCE_TESTS=true
KUBE_CONFORMANCE_IMAGE=quay.io/coreos/kube-conformance:v1.7.5_coreos.0_golang1.9.1
COMPONENT_TEST_IMAGES=quay.io/coreos/tectonic-console-tester:v2.7.1,
```

And depending on your platform:

#### AWS
``` bash
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
TF_VAR_tectonic_aws_ssh_key
TF_VAR_tectonic_aws_region
```

#### Azure
``` bash
ARM_CLIENT_ID
ARM_CLIENT_SECRET
ARM_ENVIRONMENT
ARM_SUBSCRIPTION_ID
ARM_TENANT_ID
TF_VAR_tectonic_azure_location
```

#### GCP
``` bash
GOOGLE_APPLICATION_CREDENTIALS
GOOGLE_CREDENTIALS
GOOGLE_CLOUD_KEYFILE_JSON
GCLOUD_KEYFILE_JSON
GOOGLE_PROJECT
TF_VAR_tectonic_gcp_ssh_key
```

### 2. Launch the tests
Once the environment variables are set, run `make tests/smoke
TEST=aws/basic_spec.rb`, where `TEST=<xx>` represents the test spec you want to
run.


## Running Conformance tests in a running cluster

There are two ways to run the conformance tests in an existing k8s cluster.
One way is to use the docker image that we provide, you can build using the [docker file](../images/kubernetes-e2e)
or use the existing image which is available in [quay.io](https://quay.io/repository/coreos/kube-conformance?tab=tags0).

The second approach is to use the same process as described in the [CNCF K8s certification](https://github.com/cncf/k8s-conformance/blob/master/instructions.md)

We will describe both approaches.

### 1. Running Conformance tests using `quay.io/coreos/kube-conformance`

To run the conformance testing using the docker image we build you will need to the following:

* Have a running cluster
* Kube Config file

Then, you need to execute the following command in your shell:

```Bash
$ docker run -v <KUBECONFIG_PATH>:/kubeconfig quay.io/coreos/kube-conformance:<TAG>
```

where:
`<TAG>` is the version of the Kubernetes we are using, i.e, `v1.8.2_coreos.0`
The first part of the tag `v1.8.2` means the kubernetes version.
The second part `_coreos.0` is related to hyperkube which you can find [here](https://quay.io/repository/coreos/hyperkube?tag=latest&tab=tags)
Usually the version of the conformance image should match with the cluster you have,
i.e. if you have a K8s cluster running version 1.7.5 it is better run the 1.7.5 conformance test image (in this case: `v1.7.5_coreos.0`)

The CNCF agreed for the versions 1.7.X and 1.8.X you can run the same set of tests that exist for the version 1.8.X.
Above this version (i.e. > 1.9.X) there is not guideline yet, but you might need to run the conformance tests for those versions.


### 2. Running Conformance tests using the `CNCF Process`

If you want to run the conformance tests as the same way the CNCF certification process does,
you will need to do the following:

* Have a running cluster
* Kube Config file

```Bash
$ curl -L https://raw.githubusercontent.com/cncf/k8s-conformance/master/sonobuoy-conformance.yaml | kubectl apply -f -
```

Watch Sonobuoy's logs with `kubectl logs -f -n sonobuoy sonobuoy` and wait for the line `no-exit was specified, sonobuoy is now blocking`.
When that appears you can copy the results using the following command:

```Bash
$ kubectl cp sonobuoy/sonobuoy:/tmp/sonobuoy ./results
```

#### 2.1. Submitting the results to `CNCF Certification`

If you need to submit the results to the `CNCF` you can just follow the instructions described [here](https://github.com/cncf/k8s-conformance/blob/master/instructions.md).

Which is basically run the Conformance tests described in Section 2 and create some files and then open a Pull Request in the [CNCF K8s Conformance](https://github.com/cncf/k8s-conformance).

