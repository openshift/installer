# terraform-provider-time

Terraform Provider for time-based resources.

Please note: Issues on this repository are intended to be related to bugs or feature requests with this particular provider codebase. See [Terraform Community](https://www.terraform.io/community.html) for a list of resources to ask questions about Terraform or this provider and its usage.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+

## Using the Provider

This Terraform Provider is available to install automatically via `terraform init`. It is recommended to setup the following Terraform configuration to pin the major version:

```hcl
# Terraform 0.13 and later
terraform {
  required_providers = {
    time = {
      source  = "hashicorp/time"
      version = "~> X.Y" # where X.Y is the current major version and minor version
    }
  }
}

# Terraform 0.12
terraform {
  required_providers = {
    time = "~> X.Y" # where X.Y is the current major version and minor version
  }
}
```

Additional documentation, including available resources and their arguments/attributes can be found on the [Terraform documentation website](https://terraform.io/docs/providers/time).

## Developing the Provider

If you wish to work on the provider, you'll first need [Go 1.16 or later](http://www.golang.org) installed on your machine. This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH).

### Building the Provider

From the top directory of the repository:

```console
$ go build
```

A `terraform-provider-time` binary will be left in the current directory.

### Running the Custom Provider

Follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins), e.g.place the custom provider into your plugins directory and run `terraform init` to initialize it.

### Testing the Provider

From the top directory of the repository:

```console
$ go test ./...
```
