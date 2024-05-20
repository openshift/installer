# oVirt Terraform provider

This is a Terraform provider for [oVirt](https://ovirt.org).

## Using this provider

This provider can be used with Terraform 0.13+ from the Terraform registry:

```hcl
terraform {
  required_providers {
    ovirt = {
      source = "ovirt/ovirt"
    }
  }
}

provider "ovirt" {
  # Configuration options
}
```

## Documentation

The detailed documentation can be found in the [Terraform registry](https://registry.terraform.io/providers/ovirt/ovirt/latest/docs).

## Contributing

If you wish to contribute to this Terraform provider please head over to the [contributing guide](CONTRIBUTING.md) for a detailed tutorial on how to write code here.

## History

The original roots of this provider were developed by [Maigard](https://github.com/Maigard) at [EMSL-MSC](http://github.com/EMSL-MSC/terraform-provider-ovirt) and was extensively worked on by [imjoey](https://github.com/imjoey). In 2021 a [complete rewrite](https://blogs.ovirt.org/2021/10/important-changes-to-the-ovirt-terraform-provider/) was made on the basis of [go-ovirt-client](https://github.com/ovirt/go-ovirt-client) to support better testability. The original provider can still be found in the [v0 branch](https://github.com/oVirt/terraform-provider-ovirt/tree/v0).
