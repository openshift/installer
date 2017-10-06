# Tectonic Installer Tests


## Running basic tests on PRs

Our basic set of tests includes:
- Code linting
- UI tests
- Backend unit tests

They are run on **every** PR by default.

### Actions required
- **none**

### FAQ
- *How do I retrigger the tests?*

  Comment with `ok to test` on the PR.


## Running smoke tests on PRs

In addition to our basic set of tests we have smoke tests. These test the
Tectonic installer on our supported platforms.
- AWS
- Azure
- Bare metal

### Actions required
- Add the `run-smoke-tests` GitHub label
- Add the `platform/<xxx>` GitHub label for **each** platform you want to run
  tests against
  
### FAQ
- *How do I retrigger the tests?*

  Comment with `ok to test` on the PR.

- *I forgot to add the GitHub labels. Can I add them after creating the PR?*

  Yes, just add the GitHub labels and comment `ok to test` on the PR.

- *What if I only add the `run-smoke-tests` GitHub label, but no
  `platform/<xxx>` label?*

  No smoke tests will be executed.
  
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


## Running smoke tests locally

### 1. Expose environment variables
To run a smoke test locally you need to set the following environment variables:
```
CLUSTER
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
TF_VAR_tectonic_aws_region
TF_VAR_tectonic_license_path
TF_VAR_tectonic_pull_secret_path
TF_VAR_base_domain
```

> Make sure both the *Tectonic pull secret* as well as the *Tectonic license* is
> saved somewhere in the repository folder. Only the repository folder will be
> mounted into the Docker container where the tests will be executed in. The test
> framework will not be able to read any files outside the repository folder
> during test execution.

### 2. Launch the tests
Once the environment variables are set, run `make tests/smoke TEST=aws_spec.rb`.
