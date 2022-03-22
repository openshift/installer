# Terraform Nutanix Provider

Terraform provider plugin to integrate with Nutanix Enterprise Cloud

NOTE: The latest version of the Nutanix provider is [v1.2.0](https://github.com/nutanix/terraform-provider-nutanix/releases/tag/v1.2.0)

## Build, Quality Status

 [![Go Report Card](https://goreportcard.com/badge/github.com/nutanix/terraform-provider-nutanix)](https://goreportcard.com/report/github.com/nutanix/terraform-provider-nutanix)
<!-- [![Maintainability](https://api.codeclimate.com/v1/badges/8b9e61df450276bbdbdb/maintainability)](https://codeclimate.com/github/nutanix/terraform-provider-nutanix/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/8b9e61df450276bbdbdb/test_coverage)](https://codeclimate.com/github/nutanix/terraform-provider-nutanix/test_coverage) -->

| Master                                                                                                                                                          | Develop                                                                                                                                                           |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [![Build Status](https://travis-ci.org/nutanix/terraform-provider-nutanix.svg?branch=master)](https://travis-ci.org/nutanix/terraform-provider-nutanix) | [![Build Status](https://travis-ci.org/nutanix/terraform-provider-nutanix.svg?branch=develop)](https://travis-ci.org/nutanix/terraform-provider-nutanix) |

## Community

Nutanix is taking an inclusive approach to developing this new feature and welcomes customer feedback. Please see our development project on GitHub (you're here!), comment on requirements, design, code, and/or feel free to join us on Slack. Instructions on commenting, contributing, and joining our community Slack channel are all located within our GitHub Readme.

For a slack invite, please contact terraform@nutanix.com from your business email address, and we'll add you.

### Provider Development
* [Terraform](https://www.terraform.io/downloads.html) 0.12+
* [Go](https://golang.org/doc/install) 1.12+ (to build the provider plugin)

### Provider Use

The Terraform Nutanix provider is designed to work with Nutanix Prism Central, such that you can manage one or more Prism Element clusters at scale. AOS/PC 5.6.0 or higher is required, as this Provider makes exclusive use of the v3 APIs

> For the 1.2.0 release of the provider it will have an N-1 compatibility with the Prism Central APIs. This provider was tested against Prism Central versions 2020.9 and 2020.11, as well as AOS version 5.18 and 5.19

## Example Usage

See the Examples folder for a handful of main.tf demos as well as some pre-compiled binaries.

We'll be refreshing these examples and binaries as we work through tech preview.

Long term, once this is upstream, no pre-compiled binaries will be needed, as terraform will automatically download on use.

## Configuration Reference

The following keys can be used to configure the provider.

* **endpoint** - (Required) IP address for the Nutanix Prism Central.
* **username** - (Required) Username for Nutanix Prism Central. Could be local cluster auth (e.g. `auth`) or directory auth.
* **password** - (Required) Password for the provided username.
* **port** - (Optional) Port for the Nutanix Prism Central. Default port is 9440.
* **insecure** - (Optional) Explicitly allow the provider to perform insecure SSL requests. If omitted, default value is false.
* **wait_timeout** - (optional) Set if you know that the creation o update of a resource may take long time (minutes).

```hcl
provider "nutanix" {
  username     = "admin"
  password     = "myPassword"
  port         = 9440
  endpoint     = "10.36.7.201"
  insecure     = true
  wait_timeout = 10
}
```

## Resources

* nutanix_access_control_policy
* nutanix_category_key
* nutanix_category_value
* nutanix_image
* nutanix_karbon_cluster
* nutanix_karbon_private_registry
* nutanix_network_security_rule
* nutanix_project
* nutanix_protection_rule
* nutanix_recovery_plan
* nutanix_role
* nutanix_subnet
* nutanix_user
* nutanix_virtual_machine


## Data Sources

* nutanix_access_control_policies
* nutanix_access_control_policy
* nutanix_category_key
* nutanix_cluster
* nutanix_clusters
* nutanix_host
* nutanix_hosts
* nutanix_image
* nutanix_karbon_cluster_kubeconfig
* nutanix_karbon_cluster_ssh
* nutanix_karbon_cluster
* nutanix_karbon_clusters
* nutanix_karbon_private_registries
* nutanix_karbon_private_registry
* nutanix_network_security_rule
* nutanix_permission
* nutanix_permissions
* nutanix_project
* nutanix_projects
* nutanix_role
* nutanix_roles
* nutanix_subnet
* nutanix_subnets
* nutanix_user_group
* nutanix_user_groups
* nutanix_user
* nutanix_users
* nutanix_virtual_machine
* nutanix_protection_rule
* nutanix_protection_rules
* nutanix_recovery_plan
* nutanix_recovery_plans


## Quick Install

### Install Dependencies

* [Terraform](https://www.terraform.io/downloads.html) 0.12+

### For developing or build from source


* [Go](https://golang.org/doc/install) 1.12+ (to build the provider plugin)


### Building/Developing Provider

We recomment to use Go 1.12+ to be able to use `go modules`

```sh
$ git clone https://github.com/nutanix/terraform-provider-nutanix.git
```

Enter the provider directory and build the provider

```sh
$ make tools
$ make build
```

This will create a binary file `terraform-provider-nutanix` you can copy to your terraform specific project.

Alternative build: with our demo

```sh
$ make tools
$ go build -o examples/terraform-provider-nutanix
$ cd examples
$ terraform init #to try out our demo
```

If you need multi-OS binaries such as Linux, macOS, Windows. Run the following command.

```sh
$ make tools
$ make cibuild
```

This coommand will create a `pkg/` directory with all the binaries for the most popular OS.


### Common Issues using the development binary.

Terraform download the released binary instead developent one.

Just follow this steps to get the development binary:

1. Copy the development terraform binary in the root folder of the project (i.e. where your main.tf is), this should be named `terraform-provider-nutanix`
2. Remove the entire “.terraform” directory.
    ```sh
    rm -rf .terraform/
    ```

3. Run the following command in the same folder where you have copied the development terraform binary.
    ```sh
    terraform init -upgrade
    terraform providers -version
    ```

4. You should see version as “nutanix (unversioned)”
5. Then run your main.tf

## Release it

1. Install `goreleaser` tool:

    ```bash
    go get -v github.com/goreleaser/goreleaser
    cd $GOPATH/src/github.com/goreleaser/goreleaser
    go install
    ```

    Alternatively you can download a latest release from [goreleaser Releases Page](https://github.com/goreleaser/goreleaser/releases)

1. Clean up folder `(builds)` if exists

1. Make sure that the repository state is clean:

    ```bash
    git status
    ```

1. Tag the release:

    ```bash
    git tag v1.1.0
    ```

1. Run `goreleaser`:

    ```bash
    cd (TODO: go dir)
    goreleaser --skip-publish v1.1.0
    ```

1. Check builds inside `(TODO: build dir)` directory.

1. Publish release tag to GitHub:

    ```bash
    git push origin v1.1.0
    ```

## Additional Resources

We've got a handful of resources outside of this repository that will help users understand the interactions between terraform and Nutanix

* YouTube
  _ Overview Video: [](https://www.youtube.com/watch?v=V8_Lu1mxV6g)
  _ Working with images: [](https://www.youtube.com/watch?v=IW0eQevZ73I)
* Nutanix GitHub
  _ [](https://github.com/nutanix/terraform-provider-nutanix)
  _ Private repo until code goes upstream
* Jon’s GitHub
  _ [](https://github.com/JonKohler/ThisOldCloud/tree/master/Terraform-Nutanix)
  _ Contains sample TF’s and PDFs from the youtube videos
* Slack channel \* User community slack channel is available on nutanix.slack.com. Email terraform@nutanix.com to gain entry.
