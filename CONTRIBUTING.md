# How to Contribute

CoreOS projects are [Apache 2.0 licensed](LICENSE) and accept contributions via
GitHub pull requests. This document outlines some of the conventions on
development workflow, commit message formatting, contact points and other
resources to make it easier to get your contribution accepted.

# Tectonic Installer contributions

Tectonic Installer provides specific guidelines for the modification of included Terraform modules. For more information, please see [Modifying Tectonic Installer][modify-installer].

For more information on Terraform, please see the [Terraform Documentation][tf-doc].

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](DCO) file for details.

##  Reporting a security vulnerability

Due to their public nature, GitHub and mailing lists are not appropriate places for reporting vulnerabilities. Please refer to CoreOS's [security disclosure][disclosure] process when reporting issues that may be security related.

## Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for build and test instructions
- Play with the project, submit bugs, submit patches!

### Contribution Flow

Anyone may [file issues][new-issue].
For contributors who want to work up pull requests, the workflow is roughly:

1. Create a topic branch from where you want to base your work (usually master).
2. Make commits of logical units.
3. Make sure your commit messages are in the proper format (see [below](#commit-message-format)).
4. Push your changes to a topic branch in your fork of the repository.
5. Make sure the tests pass, and add any new tests as appropriate.
6. Please run this command before submitting your pull request:
    ```sh
    make structure-check
    ```
    Note that a portion of the docs and examples are generated and that the generated files are to be committed by you. `make structure-check` checks that what is generated is what you must commit.
7. Submit a pull request to the original repository.
8. The [repo owners](OWNERS) will respond to your issue promptly, following [the ususal Prow workflow][prow-review].

Thanks for your contributions!

## Coding Style

The coding style suggested by the Golang community is used in installer. See the [style doc][golang-style] for details. Please follow them when working on your contributions.

Tectonic Installer includes syntax checks on the Terraform templates which will fail the PR checker for non-standard formatted code.

Use `make structure-check` to identify files that don't meet the canonical format and style. Then, use `terraform fmt` to align the template syntax, if necessary.

## Commit Message Format

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

```
scripts: add the test-cluster command

this uses tmux to set up a test cluster that you can easily kill and
start for debugging.

Fixes #38
```

The format can be described more formally as follows:

```
<subsystem>: <what changed>
<BLANK LINE>
<why this change was made>
<BLANK LINE>
<footer>
```

The first line is the subject and should be no longer than 70 characters, the
second line is always blank, and other lines should be wrapped at 80 characters.
This allows the message to be easier to read on GitHub as well as in various
git tools.


[modify-installer]: Documentation/contrib/modify-installer.md
[tf-doc]: https://www.terraform.io/docs/index.html
[golang-style]: https://github.com/golang/go/wiki/CodeReviewComments
[disclosure]: https://coreos.com/security/disclosure/
[new-issue]: https://github.com/openshift/installer/issues/new
[prow-review]: https://github.com/kubernetes/community/blob/master/contributors/guide/owners.md#the-code-review-process
