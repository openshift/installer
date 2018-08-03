# Tectonic Installer Tests


## Running basic tests on PRs

Our basic set of tests includes:
- Code linting
- Backend unit tests

They are run on **every** PR by default. Successful basic tests are required in
order to merge any PRs. Before starting to test your proposed changes, they are
temporarily merged into the target branch of the pull request.

### Actions required
- **none**


## Running smoke

In addition to our basic set of tests we have smoke tests which are running on AWS platform only.

### Actions required
- Add the `run-smoke-tests` GitHub label

### FAQ
- *I am not able to add labels, what should I do?*

  Please ask one of [the repository maintainers](../OWNERS) to add the
  labels.

- *How do I retrigger the tests?*

  comment with `ok to test` on the PR.

- *I forgot to add the GitHub labels. Can I add them after creating the PR?*

  Yes, just add the GitHub labels and comment `ok to test` on the PR.

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

We use Jenkins to run our tests and our `nightly` builds against the master branch.
Those jobs are executing smoke tests and building docker images.

## Running smoke tests locally

### 1. Expose environment variables

To run the smoke tests locally you need to set the following
environment variables:
``` bash
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
LICENSE_PATH
PULL_SECRET_PATH
```

Optionally you can also set:
```bash
DOMAIN
AWS_REGION
```

### 2. Launch the tests
Once the environment variables are set, run `./tests/run.sh`.
