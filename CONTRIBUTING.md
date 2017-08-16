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

## Email and Chat

The project currently uses the general CoreOS email list and IRC channel:
- Email: [coreos-dev](https://groups.google.com/forum/#!forum/coreos-dev)
- IRC: #[coreos](irc://irc.freenode.org:6667/#coreos) IRC channel on freenode.org

Please avoid emailing maintainers found in the MAINTAINERS file directly. They
are very busy and read the mailing lists.

##  Reporting a security vulnerability

Due to their public nature, GitHub and mailing lists are not appropriate places for reporting vulnerabilities. Please refer to CoreOS's [security disclosure][disclosure] process when reporting issues that may be security related.

## Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for build and test instructions
- Play with the project, submit bugs, submit patches!
- Go get a couple of projects nesessary for updating docs and examples:
    ```shell
    $ go get github.com/segmentio/terraform-docs
    $ go get github.com/s-urbaniak/terraform-examples
    ```

### Contribution Flow

This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where you want to base your work (usually master).
- Make commits of logical units.
- Make sure your commit messages are in the proper format (see below).
- Push your changes to a topic branch in your fork of the repository.
- Make sure the tests pass, and add any new tests as appropriate.
- Please run this command before submitting your pull request:
    ```shell
    $ make structure-check
    ```
- Note that a portion of the docs and examples are generated and that the generated files are to be committed by you. `make structure-check` checks that what is generated is what you must commit.
- Submit a pull request to the original repository.

Thanks for your contributions!

## Coding Style

CoreOS projects written in Go follow a set of style guidelines that we've documented in [Go at CoreOS][coreos-golang]. Please follow them when
working on your contributions.


Tectonic Installer includes syntax checks on the Terraform templates which will fail the PR checker for non-standard formatted code.

Use `make structure-check` to identify files that don't meet the canonical format and style. Then, use `terraform fmt` to align the template syntax, if necessary.

## Format of the Commit Message

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
[coreos-golang]: https://github.com/coreos/docs/tree/master/golang
[disclosure]: https://coreos.com/security/disclosure/