<!-- markdownlint-disable first-line-h1 no-inline-html -->
<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for VMware vSphere

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/hashicorp/terraform-provider-vsphere?label=release&style=for-the-badge)](https://github.com/hashicorp/terraform-provider-vsphere/releases/latest) [![License](https://img.shields.io/github/license/hashicorp/terraform-provider-vsphere.svg?style=for-the-badge)](LICENSE)

The Terraform Provider for VMware vSphere is a plugin for Terraform that allows you to interact with VMware vSphere, notably [vCenter Server][vmware-vcenter] and [ESXi][vmware-esxi]. This provider can be used to manage a VMware vSphere environment, including virtual machines, host and cluster management, inventory, networking, storage, datastores, content libraries, and more.

Learn more:

* Read the provider [documentation][provider-documentation].

* Join the community [discussions][provider-discussions].

## Requirements

* [Terraform 0.13+][terraform-install]

    For general information about Terraform, visit [terraform.io][terraform-install] and [the project][terraform-github] on GitHub.

* [Go 1.19][golang-install]

    Required if building the provider.

* [VMware vSphere][vmware-vsphere-documenation]

    The provider supports VMware vSphere versions in accordance with the VMware Product Lifecycle Matrix from General Availability to End of General Support.

    Learn more: [VMware Product Lifecycle Matrix][vmware-product-lifecycle-matrix]

    > **NOTE**
    >
    > This provider requires API write access and is therefore **not supported** for use with a free VMware vSphere Hypervisor license.

## Using the Provider

The Terraform Provider for VMware vSphere is an official provider. Official providers are maintained by the Terraform team at [HashiCorp][hashicorp] and are listed on the [Terraform Registry][terraform-registry].  

To use a released version of the Terraform provider in your environment, run `terraform init` and Terraform will automatically install the provider from the Terraform Registry.

Unless you are contributing to the provider or require a pre-release bugfix or feature, use an **officially** released version of the provider.

See [Installing the Terraform Provider for VMware vSphere][provider-install] for additional instructions on automated and manual installation methods and how to control the provider version.

For either installation method, documentation about the provider configuration, resources, and data sources can be found on the Terraform Registry.

## Upgrading the Provider

The provider does not upgrade automatically. After each new release, you can run the following command to upgrade the provider:

```shell
terraform init -upgrade
```

## Contributing

The Terraform Provider for VMware vSphere is the work of many contributors and the project team appreciates your help!

If you discover a bug or would like to suggest an enhancement, submit [an issue][provider-issues]. Once submitted, your issue will follow the [lifecycle][provider-issue-lifecycke] process.

If you would like to submit a pull request, please read the [contribution guidelines][provider-contributing] to get started. In case of enhancement or feature contribution, we kindly ask you to open an issue to discuss it beforehand.

Learn more in the [Frequently Asked Questions][provider-faq].

## License

The Terraform Provider for VMware vSphere is available under the [Mozilla Public License, version 2.0][provider-license] license.

[golang-install]: https://golang.org/doc/install
[hashicorp]: https://hashicorp.com
[provider-contributing]: docs/CONTRIBUTING.md
[provider-documentation]: https://registry.terraform.io/providers/hashicorp/vsphere/latest/docs
[provider-discussions]: https://discuss.hashicorp.com/tags/c/terraform-providers/31/vsphere
[provider-faq]: docs/FAQ.md
[provider-install]: docs/INSTALL.md
[provider-issues]: https://github.com/hashicorp/terraform-provider-vsphere/issues/new/choose
[provider-issue-lifecycke]: docs/ISSUES.md
[provider-license]: LICENSE
[terraform-install]: https://www.terraform.io/downloads.html
[terraform-github]: https://github.com/hashicorp/terraform
[terraform-registry]: https://registry.terraform.io
[vmware-esxi]: https://www.vmware.com/products/esxi-and-esx.html
[vmware-product-lifecycle-matrix]: https://lifecycle.vmware.com
[vmware-vcenter]: https://www.vmware.com/products/vcenter-server.html
[vmware-vsphere-documenation]: https://docs.vmware.com/en/VMware-vSphere
