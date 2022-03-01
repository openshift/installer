# vSphere Provider for Terraform [![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/hashicorp/terraform-provider-vsphere?label=release)](https://github.com/hashicorp/terraform-provider-vsphere/releases) [![license](https://img.shields.io/github/license/hashicorp/terraform-provider-vsphere.svg)]()


<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terrafpr," align="right" height="50" />
</a>

* [Getting Started & Documentation](https://www.terraform.io/docs/providers/vsphere/index.html)
* Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)


This is the repository for the vSphere Provider for Terraform, which one can use
with Terraform to work with VMware vSphere Products, notably [vCenter
Server][vmware-vcenter] and [ESXi][vmware-esxi].

[vmware-vcenter]: https://www.vmware.com/products/vcenter-server.html
[vmware-esxi]: https://www.vmware.com/products/esxi-and-esx.html

For general information about Terraform, visit the [official
website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform

This provider plugin is maintained by the Terraform team at [HashiCorp](https://www.hashicorp.com/).

## Requirements
-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
    - Note that version 0.11.x currently works, but is [deprecated](https://www.hashicorp.com/blog/deprecating-terraform-0-11-support-in-terraform-providers/)
- vSphere 6.5    
   -  Currently, this provider is not tested for vSphere 7, but plans are underway to add support.
-	[Go](https://golang.org/doc/install) 1.14.x (to build the provider plugin)

## Building The Provider

Unless you are [contributing](_about/CONTRIBUTING.md) to the provider or require a
pre-release bugfix or feature, you will want to use an [officially released](https://github.com/hashicorp/terraform-provider-vsphere/releases)
version of the provider.


## Contributing to the provider

The vSphere Provider for Terraform is the work of many contributors. We appreciate your help!

### Trending contributors

[![](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/images/0)](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/links/0)[![](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/images/1)](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/links/1)[![](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/images/2)](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/links/2)[![](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/images/3)](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/links/3)[![](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/images/4)](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/links/4)[![](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/images/5)](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/links/5)[![](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/images/6)](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/links/6)[![](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/images/7)](https://sourcerer.io/fame/bill-rich-private/hashicorp/terraform-provider-vsphere/links/7)

To contribute, please read the [contribution guidelines](_about/CONTRIBUTING.md). You may also [report an issue](https://github.com/hashicorp/terraform-provider-vsphere/issues/new/choose). Once you've filed an issue, it will follow the [issue lifecycle](_about/ISSUES.md).

Also available are some answers to [Frequently Asked Questions](_about/FAQ.md).


